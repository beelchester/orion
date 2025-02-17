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

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Download(url string) string {
	fmt.Println("call")
	c, err := arigo.Dial("ws://localhost:6800/jsonrpc", "")
	if err != nil {
		panic(err)
	}

	// Download is a utility method which calls AddURI and
	// then waits for the download to complete
	status, err := c.Download(arigo.URIs(url), nil)
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
	return status.GID + " downloaded"
}
