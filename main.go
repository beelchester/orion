package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
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

func startAria2c() (*exec.Cmd, error) {
	running, err := isProcessRunning("aria2c")
	if err != nil {
		return nil, err
	}
	if running {
		err := fmt.Errorf("aria2c is already running")
		return nil, err
	} else {
		cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-all")
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
	aria2cCmd, err := startAria2c()
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

