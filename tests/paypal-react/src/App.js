import React, { useEffect, useState } from "react";
import { PayPalScriptProvider, PayPalButtons } from "@paypal/react-paypal-js";

// Renders errors or successfull transactions on the screen.
function Message({ content }) {
  return <p>{content}</p>;
}

function App() {
  const baseOptions = {
    "client-id": undefined,
    "enable-funding": "venmo",
    "disable-funding": "",
    currency: "USD",
    "data-page-type": "product-details",
    components: "buttons",
    "data-sdk-integration-source": "developer-studio",
  }
  const [initialOptions, setInitialOptions] = useState(baseOptions);
  const [message, setMessage] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/api/v1/payment/config/paypal")
      .then(response => {
        return response.json();
      }).then(data => {
        setInitialOptions({
          ...baseOptions,
          "client-id": data.paypalClientId,
        })
      }).catch(error => {
        console.log(error);
      });

      console.log('initialOptions');
      console.log(initialOptions);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const PayPalButton = () => {
    return (
      <PayPalScriptProvider options={initialOptions}>
        <PayPalButtons
          style={{
            shape: "rect",
            layout: "vertical",
            color: "gold",
            label: "paypal",
          }}
          createOrder={async () => {
            try {
              const response = await fetch("http://localhost:8080/api/v1/payment/create-payment-intent/paypal", {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                // use the "body" param to optionally pass additional order information
                // like product ids and quantities
                body: JSON.stringify({
                  amount: 100.0,
                  currency: "USD",
                }),
              });

              const orderData = await response.json();
              if (orderData.paypalData.id) {
                return orderData.paypalData.id;
              } else {
                const errorDetail = orderData?.details?.[0];
                const errorMessage = errorDetail
                  ? `${errorDetail.issue} ${errorDetail.description} (${orderData.debug_id})`
                  : JSON.stringify(orderData);

                throw new Error(errorMessage);
              }
            } catch (error) {
              console.error(error);
              setMessage(`Could not initiate PayPal Checkout...${error}`);
            }
          }}
          onApprove={async (
            data,
            actions
          ) => {
            try {
              const response = await fetch(
                `http://localhost:8080/api/orders/${data.orderID}/capture/paypal`,
                {
                  method: "POST",
                  headers: {
                    "Content-Type": "application/json",
                  },
                }
              );

              const orderData = await response.json();
              // Three cases to handle:
              //   (1) Recoverable INSTRUMENT_DECLINED -> call actions.restart()
              //   (2) Other non-recoverable errors -> Show a failure message
              //   (3) Successful transaction -> Show confirmation or thank you message

              const errorDetail = orderData?.details?.[0];

              if (errorDetail?.issue === "INSTRUMENT_DECLINED") {
                // (1) Recoverable INSTRUMENT_DECLINED -> call actions.restart()
                // recoverable state, per https://developer.paypal.com/docs/checkout/standard/customize/handle-funding-failures/
                return actions.restart();
              } else if (errorDetail) {
                // (2) Other non-recoverable errors -> Show a failure message
                throw new Error(
                  `${errorDetail.description} (${orderData.debug_id})`
                );
              } else {
                // (3) Successful transaction -> Show confirmation or thank you message
                // Or go to another URL:  actions.redirect('thank_you.html');
                const transaction =
                  orderData.purchase_units[0].payments.captures[0];
                setMessage(
                  `Transaction ${transaction.status}: ${transaction.id}. See console for all available details`
                );
                console.log(
                  "Capture result",
                  orderData,
                  JSON.stringify(orderData, null, 2)
                );
              }
            } catch (error) {
              console.error(error);
              setMessage(
                `Sorry, your transaction could not be processed...${error}`
              );
            }
          }} 
        />
      </PayPalScriptProvider>
    )
  }

  return (
    <div className="App">
      { initialOptions["client-id"] === undefined? 
        <Message content={"Loading..."} /> :
        PayPalButton()
      }
      <Message content={message} />
    </div>
  );
}

export default App; 
