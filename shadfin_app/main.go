package main

import (
	"context"
	"fmt"
	"shadfin/config"
	"shadfin/player"
	"time"
	"unsafe"

	"shadfin/bundle"

	"github.com/bep/debounce"
	"github.com/brys0/wui-layered/v2"
	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	// resource compilation tool assigns app.ico ID of 3
	// rsrc -manifest app.manifest -ico app.ico -o rsrc.syso
	AppIconID = 3
)

func main() {
	// Create an instance of the app structure
	windowPlayer := wui.NewWindow()
	windowPlayer.SetBounds(50, 50, 1024, 768)
	windowPlayer.SetHasBorder(true)
	windowPlayer.SetAlpha(255)
	windowPlayer.SetBackground(wui.Color(0))
	windowPlayer.SetResizable(true)
	windowPlayer.SetClassName("shadfin_app")
	windowPlayer.SetTitle("Shadfin")
	icon, err := wui.NewIconFromExeResource(AppIconID)
	if err != nil {
		panic(err)
	}
	windowPlayer.SetIcon(icon)
	// Create an instance of the app structure

	app := NewApp()
	player := player.NewPlayer()
	config := config.NewConfig()

	app.OnStartup = func() {
		println("on app start")

		// Immersive dark mode
		var t int32
		t = 1
		w32.DwmSetWindowAttribute(w32.HWND(windowPlayer.Handle()), 20, w32.LPCVOID(unsafe.Pointer(&t)), uint32(unsafe.Sizeof(t)))
	}

	windowPlayer.SetOnShow(func() {
		if app.Ready {
			app.parent = windowPlayer
			player.Window = uint64(windowPlayer.Handle())
			runtime.WindowCreate(app.ctx, windowPlayer.Handle())

			println("mpv should be running handle:", windowPlayer.Handle())
		}
	})

	debounce := debounce.New(150 * time.Millisecond)
	windowPlayer.SetOnResize(func() {
		_, _, w, h := windowPlayer.InnerBounds()
		println("Main window was resized!")
		if app.Started {
			runtime.EventsEmit(app.ctx, "APP_RESIZE")
			println("Sent resize event to frontend")
			debounce(func() {
				runtime.WindowSetSize(app.ctx, w, h)
				runtime.WindowSetPosition(app.ctx, 0, 0)
				fmt.Printf("Updating webview size to: %vx%v\n", w, h)
			})
		}
	})

	windowPlayer.SetOnClose(func() {
		runtime.Quit(app.ctx)
	})

	go windowPlayer.Show()

	getContext := func(ctx context.Context) {
		app.setContext(ctx)
		player.SetContext(ctx)
		config.SetContext(ctx)
	}
	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Shadfin",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: bundle.Bundle,
		},
		DisableResize: true,
		WaitForParent: true,
		OnStartup:     app.startup,
		OnContext:     getContext,
		// OnDrag:           onDrag,
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Windows: &windows.Options{
			DisableFramelessWindowDecorations: true,
			WindowIsTranslucent:               true,
			WebviewIsTransparent:              true,
			LayeredWindow:                     true,
			ParentWindow:                      0,
			WindowClassName:                   "shadfin_webview",
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		Bind: []interface{}{
			app,
			player,
			config,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
