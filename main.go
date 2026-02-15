package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ncruces/zenity"
)

const (
	AppID      = "mohaa-handler"
	GameEnvVar = "OPM_PATH"
)

type Config struct {
	GamePath string `json:"game_path"`
}

var Schemes = []string{"mohaa", "mohaabt", "mohaash"}

func main() {
	// 1. Check if launched by a Protocol (Browser)
	if len(os.Args) > 1 && isScheme(os.Args[1]) {
		launchGame(os.Args[1])
		return
	}

	// 2. Interactive Menu
	fmt.Println("--- MOHAA URI Scheme Manager ---")
	fmt.Println("1. Install URI Schemes (mohaa, mohaabt, mohaash)")
	fmt.Println("2. Uninstall URI Schemes")
	fmt.Println("3. Set/Change Game Path (Manual)")
	fmt.Println("4. Exit")
	fmt.Print("\nSelect an option: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		install()
	case "2":
		uninstall()
	case "3":
		selectAndSavePath()
	default:
		os.Exit(0)
	}

	fmt.Println("\nProcess finished. Press Enter to close.")
	scanner.Scan()
}

func getGamePath() string {
	if p := os.Getenv(GameEnvVar); p != "" {
		return p
	}

	configDir, _ := os.UserConfigDir()
	path := filepath.Join(configDir, AppID, "config.json")
	if data, err := os.ReadFile(path); err == nil {
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err == nil && cfg.GamePath != "" {
			return cfg.GamePath
		}
	}

	return selectAndSavePath()
}

func selectAndSavePath() string {
	fmt.Println("Opening file selection dialog...")

	pattern := "*.exe"
	if runtime.GOOS != "windows" {
		pattern = "*"
	}

	path, err := zenity.SelectFile(
		zenity.Title("Select MOHAA Executable"),
		zenity.FileFilter{Name: "Executables", Patterns: []string{pattern}},
	)

	if err != nil || path == "" {
		fmt.Println("No file selected or selection cancelled.")
		return ""
	}

	configDir, _ := os.UserConfigDir()
	dir := filepath.Join(configDir, AppID)
	os.MkdirAll(dir, 0755)

	cfg := Config{GamePath: path}
	data, _ := json.Marshal(cfg)
	os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)

	fmt.Printf("Path saved successfully: %s\n", path)
	return path
}

func isScheme(arg string) bool {
	for _, s := range Schemes {
		if strings.HasPrefix(strings.ToLower(arg), s+"://") {
			return true
		}
	}
	return false
}

func launchGame(uri string) {
	path := getGamePath()
	if path == "" {
		zenity.Error("Game path not set. Please run the app directly to configure it.", zenity.Title("Error"))
		return
	}

	// Extract scheme and params
	parts := strings.SplitN(uri, "://", 2)
	if len(parts) < 2 {
		return
	}
	scheme := strings.ToLower(parts[0])
	params := strings.TrimRight(parts[1], "/")

	args := []string{"+connect", params}

	if scheme == "mohaash" {
		args = append(args, "+set", "com_target_game", "1")
	} else if scheme == "mohaabt" {
		args = append(args, "+set", "com_target_game", "2")
	}

	fmt.Printf("Launching: %s %s\n", path, strings.Join(args, " "))

	// Prepare the command
	cmd := exec.Command(path, args...)

	// --- THE FIX ---
	// Set the working directory to the folder where the .exe sits.
	// This ensures the game finds the /main folder correctly.
	gameDir := filepath.Dir(path)
	cmd.Dir = gameDir

	// Inherit system IO for logging/debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start the process
	err := cmd.Start()
	if err != nil {
		zenity.Error(fmt.Sprintf("Failed to launch game binary:\n%v", err), zenity.Title("Launch Error"))
	}
}

func install() {
	execPath, _ := os.Executable()
	fmt.Printf("Registering handler to: %s\n", execPath)

	switch runtime.GOOS {
	case "windows":
		for _, s := range Schemes {
			runCmd("reg", "add", "HKCU\\Software\\Classes\\"+s, "/ve", "/d", "URL:"+s+" Protocol", "/f")
			runCmd("reg", "add", "HKCU\\Software\\Classes\\"+s, "/v", "URL Protocol", "/d", "", "/f")
			runCmd("reg", "add", "HKCU\\Software\\Classes\\"+s+"\\shell\\open\\command", "/ve", "/d", fmt.Sprintf("\"%s\" \"%%1\"", execPath), "/f")
		}
	case "linux":
		desktopPath := filepath.Join(os.Getenv("HOME"), ".local/share/applications", AppID+".desktop")
		mimeTypeStr := "x-scheme-handler/" + strings.Join(Schemes, ";x-scheme-handler/") + ";"
		content := fmt.Sprintf("[Desktop Entry]\nName=MOHAA Launcher\nExec=\"%s\" %%u\nType=Application\nTerminal=false\nMimeType=%s", execPath, mimeTypeStr)
		os.WriteFile(desktopPath, []byte(content), 0644)
		runCmd("update-desktop-database", filepath.Dir(desktopPath))
	case "darwin":
		fmt.Println("macOS Note: To use URI schemes, ensure this binary is inside a .app bundle with a defined Info.plist.")
	}
	fmt.Println("Installation complete.")
}

func uninstall() {
	switch runtime.GOOS {
	case "windows":
		for _, s := range Schemes {
			runCmd("reg", "delete", "HKCU\\Software\\Classes\\"+s, "/f")
		}
	case "linux":
		desktopPath := filepath.Join(os.Getenv("HOME"), ".local/share/applications", AppID+".desktop")
		os.Remove(desktopPath)
		runCmd("update-desktop-database", filepath.Dir(desktopPath))
	}
	fmt.Println("Uninstalled successfully.")
}

func runCmd(name string, args ...string) {
	_ = exec.Command(name, args...).Run()
}
