package main

import (
	"context"
	"fmt"
	"strings"
	"github.com/beelchester/arigo"
)

// App struct
type App struct {
	ctx context.Context
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
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Download(url string) string {
	fmt.Println("Starting download for:", url)
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		panic(err)
	}

	// Just add the URI without waiting for completion
	// this will immediately return the GID 
	gid, err := c.AddURI([]string{url}, nil)
	
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