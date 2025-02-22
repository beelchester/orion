package main

import (
    "embed"
    "fmt"
    "log"
    "net"
    "os"
    "os/exec"
    "time"
    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

const (
    WAILS_PORT = 34115
    ARIA2C_PORT = 6800
)

func checkPort(port int) error {
    addr := fmt.Sprintf("127.0.0.1:%d", port)
    for i := 0; i < 3; i++ {
        if listener, err := net.Listen("tcp", addr); err == nil {
            listener.Close()
            return nil
        }
        exec.Command("sh", "-c", fmt.Sprintf("lsof -ti :%d | xargs kill -9", port)).Run()
        time.Sleep(time.Second)
    }
    return fmt.Errorf("port %d remains in use", port)
}

func startAria2c() (*exec.Cmd, error) {
    if err := checkPort(ARIA2C_PORT); err != nil {
        return nil, err
    }

    cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-all", "--rpc-listen-port=6800")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("aria2c start failed: %v", err)
    }
    time.Sleep(time.Second)
    return cmd, nil
}

func main() {
    checkPort(WAILS_PORT)

    aria2cCmd, err := startAria2c()
    if err != nil {
        log.Printf("Warning: aria2c start failed: %v", err)
    }

    if aria2cCmd != nil {
        defer func() {
            fmt.Println("❌ Stopping Aria2c...")
            aria2cCmd.Process.Kill()
        }()
    }

    app := NewApp()
    if err := wails.Run(&options.App{
        Title: "orion",
        Width: 1024,
        Height: 768,
        AssetServer: &assetserver.Options{Assets: assets},
        BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
        OnStartup: app.startup,
        Bind: []interface{}{app},
    }); err != nil {
        log.Printf("Error: %v", err)
    }
}