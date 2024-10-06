package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SaveCanvasImg(image string) {
	parts := strings.Split(image, ",")
	if len(parts) != 2 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Error",
			Message:       "Invalid data URI format",
			DefaultButton: "Ok",
		})
		return
	}

	data := parts[1]
	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Error",
			Message:       "Couldn't decode the image sry",
			DefaultButton: "Ok",
		})
		return
	}

	filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "canvas",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JPEG Image *.jpg",
				Pattern:     "*.jpg",
			},
		},
	})

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Error",
			Message:       "Filepath problem ;(",
			DefaultButton: "Ok",
		})
		return
	}

	if err := os.WriteFile(filepath, dec, 0644); err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "Error",
			Message:       "Couldn't create file ;(",
			DefaultButton: "Ok",
		})
	}
}
