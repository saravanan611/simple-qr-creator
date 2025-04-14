# 🧙‍♂️ QR Code Generator (Go + Tailwind + Font Awesome)

A simple, beautiful QR Code generator web app built in Go. Generate, download, and share QR codes with ease via a lightweight local web server.

## 🚀 Features

- Stylish HTML UI with TailwindCSS and Font Awesome.
- Generates PNG QR codes from any text or URL.
- Instant preview, download, and share functionality.
- Auto-opens browser to the app.
- Optional system notifications on errors.
- Self-terminates after 1 hour (customizable).

## 🛠️ Prerequisites

- **Go 1.18+**
- Optional: `notify-send` (Linux), `osascript` (macOS), or PowerShell (Windows) for desktop notifications

## 🧑‍💻 Usage

### 1. Clone the Repository

```bash
git clone https://github.com/saravanan611/simple-qr-creator.git
cd qr-generator-go
```

### 2. install module
```bash
go mod tidy
```


# You can cross-compile for any OS:

### Windows
```bash 
GOOS=windows GOARCH=amd64 go build -o app.exe main.go
```

### Linux
```bash 
GOOS=linux GOARCH=amd64 go build -o app main.go
```
### macOS
```bash 
GOOS=darwin GOARCH=amd64 go build -o app-darwin main.go
```

# Run the App

### Optional: Set a custom port
```bash 
export Qrport=5678
```

### Run the app
```bash 
./app
```