package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/loading
var loadingFS embed.FS

func main() {
	app := NewApp()

	loading, err := fs.Sub(loadingFS, "frontend/loading")
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:            "DataseAI",
		Width:            1280,
		Height:           800,
		MinWidth:         900,
		MinHeight:        600,
		DisableResize:    false,
		Fullscreen:       false,
		WindowStartState: options.Normal,
		AssetServer: &assetserver.Options{
			Assets: loading,
		},
		OnStartup: app.startup,
		OnDomReady: func(ctx context.Context) {
			runtime.WindowExecJS(ctx, fmt.Sprintf(
				`window.location.replace("http://localhost:%d")`, app.port,
			))
		},
		OnShutdown: app.shutdown,
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                 false,
				HideTitleBar:              false,
				FullSizeContent:           false,
				UseToolbar:                false,
			},
			Appearance:           mac.DefaultAppearance,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "DataseAI",
				Message: "Database management GUI",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
