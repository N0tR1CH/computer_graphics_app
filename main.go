package main

import (
	"embed"

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
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Draw Stuff",
		Width:     1280,
		Height:    720,
		MinWidth:  1280,
		MinHeight: 720,
		MaxWidth:  1280,
		MaxHeight: 0,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []any{
			app,
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
