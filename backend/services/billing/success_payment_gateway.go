package billing

import "net/http"

const (
	SuccessPaymentGatewayPath = "/success-payment-gateway"
)

// encore:api public raw method=GET path=/success-payment-gateway
func SuccessPaymentGateway(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlResponseWithSpinner))
}

// htmlResponseWithSpinner is a simple HTML page that displays a spinner and sends a message to the parent window indicating that the payment has been completed.
const htmlResponseWithSpinner = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Payment Success</title>

    <style>
        html,
        body {
            margin: 0;
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #ffffff;
            overflow: hidden;
        }

        .spinner {
            width: 48px;
            height: 48px;
            border: 4px solid #e5e7eb;
            border-top-color: #3b82f6;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }

        @keyframes spin {
            to {
                transform: rotate(360deg);
            }
        }
    </style>
</head>
<body>
    <div class="spinner"></div>

    <script>
        window.parent.postMessage(
            { type: "PAYMENT_COMPLETED" },
            "*"
        );
    </script>
</body>
</html>
`
