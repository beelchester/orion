package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// go:embed all:frontend/dist
var assets embed.FS

func isProcessRunning(name string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range processes {
		exe, err := p.Name()
		if err == nil && strings.EqualFold(exe, name) {
			return true, nil
		}
	}

	return false, nil
}

func startAria2c(downloadDir string) (*exec.Cmd, error) {
	running, err := isProcessRunning("aria2c")
	if err != nil {
		return nil, err
	}
	if running {
		err := fmt.Errorf("aria2c is already running")
		return nil, err
	} else {
		args := []string{"--enable-rpc", "--rpc-listen-all"}

		// add download dir if specified
		if downloadDir != "" {
			args = append(args, "--dir="+downloadDir)
		}

		cmd := exec.Command("aria2c", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("aria2c start failed: %v", err)
		}
		println("✅ aria2c started")
		return cmd, nil
	}
}

func main() {
	// set default download dir(user's home download folder)
	homeDir, _ := os.UserHomeDir()
	defaultDir := filepath.Join(homeDir, "Downloads")

	// if dir doesnt exist, back to home dir
	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		defaultDir = homeDir
	}

	aria2cCmd, err := startAria2c(defaultDir)
	if err != nil {
		if err.Error() == "aria2c is already running" {
			println("✅ aria2c is already running")
		} else {
			log.Printf("Error: aria2c couldn't start: %v", err)
			os.Exit(1)
		}
	}

	if aria2cCmd != nil {
		defer func() {
			// Stop aria2c, if it was started by the app
			fmt.Println("❌ Stopping Aria2c...")
			aria2cCmd.Process.Kill()
		}()
	}

	app := NewApp()
	app.downloadDir = defaultDir
	if err := wails.Run(&options.App{
		Title:            "orion",
		Width:            800,
		Height:           510,
		DisableResize:    true,
		AssetServer:      &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
	}); err != nil {
		log.Printf("Error: %v", err)
	}
}

