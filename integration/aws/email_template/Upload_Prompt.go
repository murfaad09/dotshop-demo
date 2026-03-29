package email_template

var UploadPromptTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Upload Prompt</title>
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
            background-color: #ffc107;
            color: #ffffff;
            padding: 10px 0;
            text-align: center;
        }
        .content {
            padding: 20px;
            text-align: center;
        }
        .footer {
            background-color: #ffc107;
            color: #ffffff;
            padding: 10px 0;
            text-align: center;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            margin: 20px 0;
            background-color: #ffc107;
            color: #ffffff;
            text-decoration: none;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Upload Reminder</h1>
        </div>
        <div class="content">
            <p>Dear {{.CuratorName}},</p>
            <p>It looks like you haven't uploaded anything in the past 7 days. Keep your collection fresh and engaging by adding new content regularly!</p>
            <a href="{{.UploadLink}}" class="button">Upload Now</a>
        </div>
        <div class="footer">
            <p>&copy; 2024 Your Company. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`
