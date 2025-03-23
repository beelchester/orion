package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/beelchester/arigo"
)

// App struct
type App struct {
	ctx context.Context
	downloadDir		string
	aria2cCmd		*exec.Cmd
}

// DownloadProgressInfo contains all progress information
type DownloadProgressInfo struct {
	Progress       float64 `json:"progress"`
	DownloadSpeed  float64 `json:"downloadSpeed"`  
	TotalSize      int64   `json:"totalSize"`     
	CompletedSize  int64   `json:"completedSize"`  
	Status         string  `json:"status"`        
}

// NewApp creates a new App application struct
func NewApp() *App {
	homeDir, err := os.UserHomeDir()
	defaultDir := filepath.Join(homeDir, "Downloads")

	if err != nil {
		if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
			defaultDir = homeDir
		}
	} else {
		defaultDir, _ = os.Getwd()
	}
	return &App{
		downloadDir: defaultDir,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	cmd, err := startAria2c(a.downloadDir)

	if err != nil {
		if err.Error() == "aria2c is already running" {
			println("aria2c is already running")
		} else {
			fmt.Printf("Error: aria2c couldn't start: %v\n", err)
		}
	} else {
		a.aria2cCmd = cmd
	}
}

func(a *App) GetDownloadDirectory() string {
	return a.downloadDir
}

func (a *App) SetDownloadDirectory() (string, error) {

	//open dir selection dialog
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Download Directory",
        DefaultDirectory: a.downloadDir,
	})

	if err != nil {
		return "", err
	}
	if dir == "" {
		return a.downloadDir, nil
	}

	//update download dir
	a.downloadDir  = dir


	// restart aria2c with new directory
	if a.aria2cCmd != nil {
		fmt.Println("Restarting Aria2c with new directory...")
		a.aria2cCmd.Process.Kill()

		cmd, err := startAria2c(a.downloadDir)
		if err != nil {
			return "", fmt.Errorf("failed to restart aria2c: %v", err)
		}

		a.aria2cCmd = cmd
	}

	return a.downloadDir, nil
}

func (a *App) Download(url string) string {
	fmt.Println("Starting download for:", url)
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		panic(err)
	}

	options := &arigo.Options{
		Dir: a.downloadDir,
	}

	// Just add the URI without waiting for completion
	// this will immediately return the GID 
	gid, err := c.AddURI([]string{url}, options)
	
	if err != nil {
		errMsg := err.Error()
		start := strings.Index(errMsg, "message:")
		if start != -1 {
			start += len("message:")
			end := strings.Index(errMsg[start:], "]")
			if end != -1 {
				return strings.TrimSpace(errMsg[start : start+end])
			}
		}
		return errMsg
	}
	
	// Return GID as a string
	return gid.String()
}

func (a *App) GetDownloadProgress(gid string) (*DownloadProgressInfo, error) {
	if gid == "" {
		return nil, fmt.Errorf("invalid GID")
	}
	
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		return nil, err
	}
	
	// Get status using the provided GID
	status, err := c.TellStatus(gid)
	if err != nil {
		return nil, err
	}
	
	// Calculate progress
	progress := 0.0
	if status.TotalLength > 0 {
		progress = float64(status.CompletedLength) / float64(status.TotalLength) * 100
	}
	
	// response with download info
	info := &DownloadProgressInfo{
		Progress:      progress,
		DownloadSpeed: float64(status.DownloadSpeed),
		TotalSize:     int64(status.TotalLength),
		CompletedSize: int64(status.CompletedLength),
		Status:        string(status.Status),
	}
	
	return info, nil
}

// PauseDownload pauses a download using its GID
func (a *App) PauseDownload(gid string) error {
	if gid == "" {
		return fmt.Errorf("invalid GID")
	}
	
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		return err
	}
	
	// Pause the download
	err = c.Pause(gid)
	if err != nil {
		return err
	}
	
	return nil
}

// resumes a paused download
func (a *App) ResumeDownload(gid string) error {
	if gid == "" {
		return fmt.Errorf("invalid GID")
	}
	
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		return err
	}
	
	// Unpause/resume the download
	err = c.Unpause(gid)
	if err != nil {
		return err
	}
	
	return nil
}

// cancels and removes a download
func (a *App) CancelDownload(gid string) error {
	if gid == "" {
		return fmt.Errorf("invalid GID")
	}
	
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		return err
	}
	
	err = c.Remove(gid)
	if err != nil {
		return err
	}
	
	return nil
}