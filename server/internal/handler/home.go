package handler

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/skip2/go-qrcode"
)

var homeTemplate = template.Must(template.New("home").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LaNotify</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
            color: #e0e0e0;
        }
        .container {
            background: rgba(255, 255, 255, 0.05);
            backdrop-filter: blur(10px);
            border-radius: 24px;
            padding: 40px;
            max-width: 480px;
            width: 100%;
            text-align: center;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        h1 {
            font-size: 2.5rem;
            margin-bottom: 8px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        .subtitle {
            color: #888;
            margin-bottom: 32px;
            font-size: 1.1rem;
        }
        .qr-container {
            background: white;
            border-radius: 16px;
            padding: 20px;
            display: inline-block;
            margin-bottom: 32px;
        }
        .qr-container img {
            display: block;
            width: 256px;
            height: 256px;
        }
        .instructions {
            text-align: left;
            background: rgba(255, 255, 255, 0.03);
            border-radius: 12px;
            padding: 24px;
        }
        .instructions h2 {
            font-size: 1.2rem;
            margin-bottom: 16px;
            color: #667eea;
        }
        .instructions ol {
            padding-left: 24px;
        }
        .instructions li {
            margin-bottom: 12px;
            line-height: 1.5;
        }
        .instructions li::marker {
            color: #667eea;
        }
        .note {
            margin-top: 24px;
            padding: 16px;
            background: rgba(102, 126, 234, 0.1);
            border-radius: 8px;
            font-size: 0.9rem;
            color: #aaa;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>LaNotify</h1>
        <p class="subtitle">Forward Android notifications to your desktop</p>

        <div class="qr-container">
            <img src="data:image/png;base64,{{.QRCode}}" alt="QR Code">
        </div>

        <div class="instructions">
            <h2>Setup Instructions</h2>
            <ol>
                <li>Install <strong>LaNotify</strong> app on your Android device</li>
                <li>Open the app and tap <strong>Scan QR Code</strong></li>
                <li>Point your camera at the QR code above</li>
                <li>Grant <strong>Notification Access</strong> permission</li>
                <li>Disable <strong>Battery Optimization</strong> for the app</li>
                <li>Enable the <strong>Forward Notifications</strong> switch</li>
            </ol>
        </div>

        <div class="note">
            The QR code contains the server URL and authentication token.
            Keep it private and don't share it.
        </div>
    </div>
</body>
</html>`))

// HomeHandler returns a handler that displays the home page with QR code.
func HomeHandler(serverURL, secret string) http.HandlerFunc {
	qrData := fmt.Sprintf("%s|%s", serverURL, secret)

	qr, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qr)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := homeTemplate.Execute(w, map[string]string{
			"QRCode": qrBase64,
		}); err != nil {
			log.Printf("Failed to render home page: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
