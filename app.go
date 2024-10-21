package main

import (
	"context"
	"encoding/base64"
	"image/color"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

type Cmyk struct {
	C uint8 `json:"c"`
	M uint8 `json:"m"`
	Y uint8 `json:"y"`
	K uint8 `json:"k"`
}

type Rgb struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
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

func (a *App) RgbToCmyk(r, g, b uint8) Cmyk {
	c, m, y, k := color.RGBToCMYK(r, g, b)
	return Cmyk{c, m, y, k}
}

func (a *App) CmykToRgb(c, m, y, k uint8) Rgb {
	r, g, b := color.CMYKToRGB(c, m, y, k)
	return Rgb{r, g, b}
}
