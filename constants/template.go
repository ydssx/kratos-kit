package constants

const (
	EmailRegisterTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            border: 1px solid #ddd;
        }
        .logo {
            text-align: center;
            margin-bottom: 20px;
        }
        .verification-code {
            background-color: #f0f0f0;
            padding: 10px;
            text-align: center;
            font-size: 24px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="http://www.shortplay.fun/favicon.svg" alt="logo" style="max-width: 200px;">
        </div>
        
        <h2>Verify your email to sign up for kiwishort.com</h2>
        
        <p>We have received your registration request. In order to verify your email address, We need you to enter the following verification code:</p>
        
        <div class="verification-code">
            {{.VerificationCode}}
        </div>
        
        <p>Please enter this verification code into our verification page within 12 hours. Complete email registration verification.</p>
        
        <p>If you did not initiate this request: Please ignore this email. Or contact our customer service team for assistance.</p>
        
        <p>Thank you for your trust in our services.</p>
    </div>
</body>
</html>
`

	EmailLoginTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            border: 1px solid #ddd;
        }
        .logo {
            text-align: center;
            margin-bottom: 20px;
        }
        .verification-code {
            background-color: #f0f0f0;
            padding: 10px;
            text-align: center;
            font-size: 24px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="http://www.shortplay.fun/favicon.svg" alt="logo" style="max-width: 200px;">
        </div>
        
        <h2>Verify your email to log in kiwishort.com</h2>
        
        <p>We received your login request. To verify your email address, We need you to enter the following verification code:</p>
        
        <div class="verification-code">
            {{.VerificationCode}}
        </div>
        
        <p>Please enter this verification code into our verification page within 5 minutes. Complete the email login operation.</p>
        
        <p>If you did not initiate this request: Please ignore this email. Or contact our customer service team for assistance.</p>
        
        <p>Thank you for your trust in our services.</p>
    </div>
</body>
</html>
`
)
