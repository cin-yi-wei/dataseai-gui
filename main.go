package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
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

	// The native menu is only used on macOS (it provides the app menu and
	// enables clipboard shortcuts). On Windows/Linux it renders as a light
	// Win32 menu bar that ignores the dark theme, so skip it there — WebView2
	// already handles Ctrl+C/V/X/Z/A natively.
	var appMenu *menu.Menu
	if goruntime.GOOS == "darwin" {
		appMenu = buildMenu()
	}

	err = wails.Run(&options.App{
		Title:            "DataseAI",
		Width:            1100,
		Height:           680,
		MinWidth:         760,
		MinHeight:        520,
		DisableResize:    false,
		Fullscreen:       false,
		WindowStartState: options.Normal,
		Menu:             appMenu,
		AssetServer: &assetserver.Options{
			Assets: loading,
		},
		OnStartup: app.startup,
		OnDomReady: func(ctx context.Context) {
			runtime.WindowExecJS(ctx, fmt.Sprintf(
				`if (window.location.host !== "localhost:%d") { window.location.replace("http://localhost:%d"); }`, app.port, app.port,
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
				Message: "Database management GUI\nhttps://dataseai.conray.top",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func buildMenu() *menu.Menu {
	m := menu.NewMenu()

	// App menu (macOS only — first menu becomes the app menu)
	app := m.AddSubmenu("DataseAI")
	app.AddText("About DataseAI", nil, func(_ *menu.CallbackData) {})
	app.AddSeparator()
	app.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		// runtime.Quit not available here without ctx; Wails handles Cmd+Q natively
	})

	// Edit menu — enables Cmd+C/V/X/A/Z in the webview
	edit := m.AddSubmenu("Edit")
	edit.AddText("Undo", keys.CmdOrCtrl("z"), func(_ *menu.CallbackData) {})
	edit.AddText("Redo", keys.Combo("z", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {})
	edit.AddSeparator()
	edit.AddText("Cut", keys.CmdOrCtrl("x"), func(_ *menu.CallbackData) {})
	edit.AddText("Copy", keys.CmdOrCtrl("c"), func(_ *menu.CallbackData) {})
	edit.AddText("Paste", keys.CmdOrCtrl("v"), func(_ *menu.CallbackData) {})
	edit.AddText("Select All", keys.CmdOrCtrl("a"), func(_ *menu.CallbackData) {})

	return m
}
