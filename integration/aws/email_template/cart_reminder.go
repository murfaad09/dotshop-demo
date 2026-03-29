package email_template

var CartReminderTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your Cart is Waiting</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border: 1px solid #dddddd;
        }
        .header {
            background-color: #ff5733;
            color: #ffffff;
            padding: 10px 0;
            text-align: center;
        }
        .content {
            padding: 20px;
            text-align: center;
        }
        .footer {
            background-color: #ff5733;
            color: #ffffff;
            padding: 10px 0;
            text-align: center;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            margin: 20px 0;
            background-color: #ff5733;
            color: #ffffff;
            text-decoration: none;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Your Cart is Waiting</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>You have items in your cart that are waiting for you. Complete your purchase before they sell out!</p>
            <a href="{{.CartLink}}" class="button">Check Out Now</a>
        </div>
        <div class="footer">
            <p>&copy; 2024 Your Company. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`
