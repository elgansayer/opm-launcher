# MOHAA Cross-Platform URI Handler

A lightweight utility to register and handle custom URI schemes (`mohaa://`, `mohaabt://`, `mohaash://`) for Medal of Honor: Allied Assault.

## Features
- **Cross-Platform**: Supports Windows (Registry) and Linux (XDG Desktop).
- **Single Binary**: No complex installation. Run the file to install/uninstall.
- **Param Handling**: Automatically passes IP/Port to the game using the `+connect` command.

## Setup
1. Set an environment variable `OPM_PATH` pointing to your game executable (e.g., `C:\Games\MOHAA\mohaa.exe`).
2. Run the binary and select **Option 1** to register the URI schemes.
3. Click any `mohaa://<server-ip>` link in your browser to launch the game.

## Building
Requires [Go](https://go.dev/).

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o bin/mohaa-handler.exe main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/mohaa-handler main.go