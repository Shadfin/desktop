package main

import (
	"context"
	"unsafe"

	"github.com/brys0/wui-layered/v2"
	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	parent  *wui.Window
	ctx     context.Context
	Ready   bool
	Started bool

	OnStartup func()
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.Started = true
	runtime.WindowSetPosition(a.ctx, 0, 0)

	if a.OnStartup != nil {
		a.OnStartup()
	}
}

func (a *App) setContext(ctx context.Context) {
	a.ctx = ctx
	a.Ready = true
}

func (a *App) Fullscreen() {
	a.parent.Fullscreen()
}

func (a *App) UnFullscreen() {
	a.parent.UnFullscreen()
}

func (a *App) SetBackground(r int64, g int64, b int64) {
	if a.parent == nil {
		return
	}

	color := wui.RGB(uint8(r), uint8(g), uint8(b))
	a.parent.SetBackground(color)
	unsafeSetCaptionColor(uint32(color), a.parent)
}

func unsafeSetCaptionColor(color uint32, window *wui.Window) {
	w32.DwmSetWindowAttribute(w32.HWND(window.Handle()), 35, w32.LPCVOID(unsafe.Pointer(&color)), uint32(unsafe.Sizeof(color)))
}
