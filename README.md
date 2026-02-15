# MOHAA Cross-Platform URI Handler

A lightweight utility to register and handle custom URI schemes (`mohaa://`, `mohaabt://`, `mohaash://`) for Medal of Honor: Allied Assault and its expansions.

## Features
- **Cross-Platform**: Supports Windows, Linux, and macOS.
- **Auto-Config**: Automatically detects expansions (`mohaabt`, `mohaash`) and appends the correct `com_target_game` parameter.
- **Easy Setup**: Run the binary to install URI schemes and select your game executable via a native file dialog.
- **Flexible Path**: Saves your game path in a local config file, or uses the `OPM_PATH` environment variable as a fallback.

## Setup
1. Download the latest binary for your OS from the [Releases](https://github.com/elgansayer/opm-launcher/releases) page.
2. Run the binary.
3. Select **Option 3** to set your game path (if not using `OPM_PATH`).
4. Select **Option 1** to register the URI schemes.
5. Click any `mohaa://<server-ip>` link in your browser to launch the game.

## Building
Requires [Go](https://go.dev/).

```bash
# Build for all platforms
make build-all
```
Artifacts will be located in the `bin/` directory.