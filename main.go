package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	AppID      = "mohaa-handler"
	GameEnvVar = "OPM_PATH"
)

var Schemes = []string{"mohaa", "mohaabt", "mohaash"}

func main() {
	// 1. Check if launched by a Protocol (Browser)
	if len(os.Args) > 1 {
		input := os.Args[1]
		for _, s := range Schemes {
			if strings.HasPrefix(strings.ToLower(input), s+"://") {
				launchGame(input)
				return
			}
		}
	}

	// 2. Interactive Menu for User
	fmt.Println("--- MOHAA URI Scheme Manager ---")
	fmt.Println("1. Install URI Schemes (mohaa, mohaabt, mohaash)")
	fmt.Println("2. Uninstall URI Schemes")
	fmt.Println("3. Exit")
	fmt.Print("\nSelect an option: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		install()
	case "2":
		uninstall()
	default:
		os.Exit(0)
	}

	fmt.Println("\nTask complete. Press Enter to close.")
	scanner.Scan()
}

func install() {
	execPath, _ := os.Executable()
	fmt.Printf("Registering to: %s\n", execPath)

	switch runtime.GOOS {
	case "windows":
		for _, s := range Schemes {
			runCmd("reg", "add", "HKCU\\Software\\Classes\\"+s, "/ve", "/d", "URL:"+s+" Protocol", "/f")
			runCmd("reg", "add", "HKCU\\Software\\Classes\\"+s, "/v", "URL Protocol", "/d", "", "/f")
			runCmd("reg", "add", "HKCU\\Software\\Classes\\"+s+"\\shell\\open\\command", "/ve", "/d", fmt.Sprintf("\"%s\" \"%%1\"", execPath), "/f")
		}
	case "linux":
		desktopPath := filepath.Join(os.Getenv("HOME"), ".local/share/applications", AppID+".desktop")
		content := fmt.Sprintf("[Desktop Entry]\nName=MOHAA Launcher\nExec=%s %%u\nType=Application\nTerminal=false\nMimeType=x-scheme-handler/%s;", execPath, strings.Join(Schemes, ";x-scheme-handler/"))
		os.WriteFile(desktopPath, []byte(content), 0644)
		runCmd("update-desktop-database", filepath.Dir(desktopPath))
		for _, s := range Schemes {
			runCmd("xdg-mime", "default", AppID+".desktop", "x-scheme-handler/"+s)
		}
	case "darwin":
		fmt.Println("Note: On macOS, URI schemes are usually handled via the Info.plist in the .app bundle.")
	}
	fmt.Println("Schemes installed successfully.")
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
	fmt.Println("Schemes removed.")
}

func launchGame(uri string) {
	// Strip the scheme prefix
	parts := strings.SplitN(uri, "://", 2)
	if len(parts) < 2 {
		return
	}
	// Clean the IP/Port (remove trailing slashes from browsers)
	params := strings.TrimRight(parts[1], "/")

	gamePath := os.Getenv(GameEnvVar)
	if gamePath == "" {
		fmt.Println("Error: Environment variable OPM_PATH not set.")
		// Stay open for 5 seconds so user can see error
		return
	}

	fmt.Printf("Launching %s with +connect %s\n", gamePath, params)
	cmd := exec.Command(gamePath, "+connect", params)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to launch game: %v\n", err)
	}
}

func runCmd(name string, args ...string) {
	_ = exec.Command(name, args...).Run()
}
