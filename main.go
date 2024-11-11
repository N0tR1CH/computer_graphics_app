package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/build
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create a server for image hosting
	server := &http.Server{
		Addr:    ":3000",
		Handler: http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// Create an instance of the app structure
	app := NewApp()
	worker := NewWorker(app)

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Draw Stuff",
		Width:     1280,
		Height:    720,
		MinWidth:  1280,
		MinHeight: 720,
		MaxWidth:  0,
		MaxHeight: 0,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnShutdown: func(ctx context.Context) {
			ctxShutDown, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctxShutDown); err != nil {
				log.Fatalf("Server shutdown failed: +%v", err)
			}
			log.Println("Server exited properly")
		},
		Bind: []any{
			app,
			worker,
		},
		EnumBind: []any{
			AllImageFormats,
		},
		Windows: &windows.Options{
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
			DisableWindowIcon:    true,
		},
		Mac: &mac.Options{
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
			Appearance:           mac.NSAppearanceNameDarkAqua,
			About: &mac.AboutInfo{
				Title:   "Draw Stuff",
				Message: "@ 2024-2025 Damian Richter",
				Icon:    icon,
			},
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 true,
				HideToolbarSeparator:       true,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
