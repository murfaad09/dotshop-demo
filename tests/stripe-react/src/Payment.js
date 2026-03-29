import { useEffect, useState } from "react";
import { Elements } from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";
import CheckoutForm from "./CheckoutForm";
// import { useParams } from "react-router-dom";

function Payment() {
  const [stripePromise, setStripePromise] = useState(null);
  const [clientSecret, setClientSecret] = useState("");
  const [loading, setLoading] = useState(false);
  const [amount, setAmount] = useState(0);

  useEffect(() => {
    fetch("http://localhost:8080/api/v1/payment/config/stripe").then(async (r) => {
      const {stripeClientKey} = await r.json();
      console.log(stripeClientKey)
      setStripePromise(loadStripe(`${stripeClientKey}`));
    });
    const params = new URLSearchParams(window.location.search);
    console.log(`The price is ${params.get("price")}`);
    setAmount(params.get("price"));
  }, []);

  useEffect(() => {
    if (amount <= 0) return
    const body = {
      "currency": "usd",
      "clientId": 1,
      "productId": 1,
      "amount": parseFloat(amount)
    }
    if (!loading) {
      setLoading(true)
      fetch("http://localhost:8080/api/v1/payment/create-payment-intent/stripe", {
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify(body),
      }).then(async (result) => {
        const { clientSecret } = await result.json();
        setClientSecret(`${clientSecret}`);
        setLoading(false)
      }).catch((error) => {
        setLoading(false)
        console.log(error)
      });
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [amount]);


  return (
    <div style={{}}>
      <h1>React Stripe and the Payment Element</h1>
      {clientSecret && stripePromise && (
        <Elements
          stripe={stripePromise}
          options={{ clientSecret, locale: "en" }}
        >
          <CheckoutForm />
        </Elements>
      )}
    </div>
  );
}

export default Payment;