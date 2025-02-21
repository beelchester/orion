package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS


func StartAria2c() (*exec.Cmd, error) {
	cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start aria2c: %v", err)
		return nil, err
	}

	fmt.Println("Aria2c started successfully.")
	return cmd, nil
}

func main() {
	// Start aria2c
	aria2cCmd, err := StartAria2c()
	if err != nil {
		log.Fatal(err)
	}

	// killed when the app terminate
	defer func() {
		if aria2cCmd != nil {
			fmt.Println("❌ Stopping Aria2c...")
			aria2cCmd.Process.Kill()
		}
	}()


	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "orion",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
