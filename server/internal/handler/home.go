package handler

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/skip2/go-qrcode"
)

//go:embed lanotifica.png
var logoPNG []byte

var homeTemplate = template.Must(template.New("home").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LaNotifica</title>
    <link rel="icon" type="image/png" href="/favicon.png">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: #0a0a0a;
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            color: #fafafa;
            line-height: 1.6;
        }

        .container {
            display: grid;
            grid-template-columns: auto 1fr;
            gap: 80px;
            max-width: 900px;
            padding: 60px;
            align-items: center;
        }

        .qr-section {
            display: flex;
            flex-direction: column;
            align-items: center;
        }

        .qr-wrapper {
            background: #fff;
            padding: 20px;
            border-radius: 24px;
            box-shadow: 0 0 0 1px rgba(255,255,255,0.1), 0 25px 50px -12px rgba(0,0,0,0.5);
        }

        .qr-wrapper img {
            display: block;
            width: 280px;
            height: 280px;
        }

        .qr-hint {
            margin-top: 20px;
            font-size: 13px;
            color: #666;
            text-align: center;
        }

        .content {
            max-width: 400px;
        }

        .header {
            display: flex;
            align-items: center;
            gap: 16px;
            margin-bottom: 12px;
        }

        .header img {
            width: 128px;
            height: 128px;
            border-radius: 20px;
        }

        h1 {
            font-size: 3rem;
            font-weight: 600;
            letter-spacing: -0.03em;
            background: linear-gradient(135deg, #fff 0%, #999 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .tagline {
            font-size: 1.1rem;
            color: #666;
            margin-bottom: 48px;
            font-weight: 300;
        }

        .steps {
            list-style: none;
            counter-reset: step;
        }

        .steps li {
            counter-increment: step;
            display: flex;
            align-items: flex-start;
            gap: 16px;
            margin-bottom: 20px;
            font-size: 15px;
            color: #a1a1a1;
        }

        .steps li::before {
            content: counter(step);
            flex-shrink: 0;
            width: 28px;
            height: 28px;
            background: #1a1a1a;
            border: 1px solid #333;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            font-weight: 500;
            color: #fff;
        }

        .steps strong {
            color: #fff;
            font-weight: 500;
        }

        .steps a {
            color: #60a5fa;
            text-decoration: none;
        }

        .steps a:hover {
            text-decoration: underline;
        }

        .note {
            margin-top: 40px;
            padding: 16px 20px;
            background: #111;
            border-radius: 12px;
            font-size: 13px;
            color: #555;
            border: 1px solid #222;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="qr-section">
            <div class="qr-wrapper">
                <img src="data:image/png;base64,{{.QRCode}}" alt="QR Code">
            </div>
            <p class="qr-hint">Scan with LaNotifica app</p>
        </div>

        <div class="content">
            <div class="header">
                <img src="data:image/png;base64,{{.Logo}}" alt="LaNotifica">
                <h1>LaNotifica</h1>
            </div>
            <p class="tagline">Forward Android notifications to your Linux desktop</p>

            <ol class="steps">
                <li><span>Install <strong>LaNotifica</strong> from <a href="https://play.google.com/store/apps/details?id=com.alessandrolattao.lanotifica" target="_blank">Google Play</a></span></li>
                <li><span>Open the app and tap <strong>Scan QR Code</strong></span></li>
                <li><span>Grant <strong>Notification Access</strong> permission</span></li>
                <li><span>Disable <strong>Battery Optimization</strong></span></li>
                <li><span>Enable <strong>Forward Notifications</strong></span></li>
            </ol>

            <p class="note">
                The QR contains your auth token and certificate fingerprint.
                Server discovery happens automatically via mDNS.
                <br><br>
                <span style="color: #777;">Version {{.Version}}</span>
            </p>
        </div>
    </div>
</body>
</html>`))

// FaviconHandler returns a handler that serves the favicon.
func FaviconHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write(logoPNG)
	}
}

// HomeHandler returns a handler that displays the home page with QR code.
func HomeHandler(secret, certFingerprint, version string) http.HandlerFunc {
	// QR format: token|fingerprint (URL is discovered via mDNS)
	qrData := fmt.Sprintf("%s|%s", secret, certFingerprint)

	qr, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qr)
	logoBase64 := base64.StdEncoding.EncodeToString(logoPNG)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := homeTemplate.Execute(w, map[string]string{
			"QRCode":  qrBase64,
			"Logo":    logoBase64,
			"Version": version,
		}); err != nil {
			log.Printf("Failed to render home page: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
