package player

import (
	"context"
	"fmt"
	r "runtime"
	"strconv"
	"sync"

	"github.com/gen2brain/go-mpv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	PLAYER_PAUSE            = "PLAYER_PAUSE"
	PLAYER_LOADED           = "PLAYER_LOADED"
	PLAYER_LOADING          = "PLAYER_LOADING"
	PLAYER_BUFFER           = "PLAYER_BUFFER"
	PLAYER_PLAYBACK_RESTART = "PLAYER_PLAYBACK_RESTART"
	PLAYER_POSITION         = "PLAYER_POSITION"
	PLAYER_MESSAGE          = "PLAYER_MESSAGE"
	PLAYER_GENERIC_ERROR    = "PLAYER_GENERIC_ERROR"
)

type Player struct {
	ctx context.Context
	// Current active mpv player, if any
	handle *mpv.Mpv
	// Whether player is paused
	paused bool
	// Should the event loop for mpv be running or not
	loop bool
	// Wait group for event loop
	wg sync.WaitGroup
	// Mutex lock for any events that read or write to the player handle
	mutex sync.Mutex
	// Parent window handle
	Window uint64
	// Actively playing URL
	URL string
	// Custom player options
	Options *PlayerOptions
}

func NewPlayer() *Player {
	return &Player{wg: sync.WaitGroup{}, mutex: sync.Mutex{}}
}

func (p *Player) SetContext(ctx context.Context) {
	p.ctx = ctx
}

func (p *Player) SetURL(url string) {
	p.URL = url
}

func (p *Player) Start() error {
	go p.startMPVRoutine()
	return nil
}

func (p *Player) Destroy() error {
	if p.handle == nil {
		return fmt.Errorf("player was never created thus couldn't be destroyed")
	}

	p.loop = false

	// Must call wakeup so the event loop doesn't keep running, otherwise memory corruption
	p.handle.Wakeup()
	p.handle.TerminateDestroy()

	return nil
}

// Should be started as go routine
func (p *Player) startMPVRoutine() {

	p.wg.Add(1)

	runtime.Invoke(p.ctx, func() {
		// Start another go routine for the player loop
		go p.startMPV()
	})

	p.wg.Wait()
}

func (p *Player) SetPlayerPause(paused bool) error {
	if p.handle == nil {
		return fmt.Errorf("can not pause player when no player is active")
	}

	// runtime.WindowFullscreen()
	println("Paused state ", paused)
	p.handle.SetProperty("pause", mpv.FormatFlag, paused)
	return nil
}

func (p *Player) SetPlayerPosition(position_sec float64) error {
	if p.handle == nil {
		return fmt.Errorf("can not set player position when no player is active")
	}

	println("Setting playback head to ", position_sec)

	p.handle.SetProperty("time-pos", mpv.FormatDouble, position_sec)
	return nil
}

// Should be started inside startMPVRoutine()
func (p *Player) startMPV() {
	// Make sure to not change the underlying thread context
	r.LockOSThread()

	// Once the event loop is done we can release the thread back to go
	defer r.UnlockOSThread()

	p.mutex.Lock()
	p.handle = mpv.New()

	// Stop loop and terminate mpv player
	defer func() {
		p.loop = false
		p.handle.TerminateDestroy()
	}()

	// Window options
	_ = p.handle.SetOptionString("border", "no")
	_ = p.handle.SetOptionString("window-dragging", "no")
	_ = p.handle.SetOptionString("wid", strconv.FormatUint(p.Window, 10))

	// MPV options
	_ = p.handle.RequestLogMessages("info")

	// Video out options
	_ = p.handle.SetOptionString("vo", "gpu-next")
	_ = p.handle.SetOptionString("gpu-api", "vulkan")
	_ = p.handle.SetOptionString("hwdec", "vulkan")

	// OSD Options
	_ = p.handle.SetOption("osc", mpv.FormatFlag, false)
	_ = p.handle.SetPropertyString("input-default-bindings", "no")
	_ = p.handle.SetOptionString("input-vo-keyboard", "no")

	// Watch properties
	_ = p.handle.ObserveProperty(0, "pause", mpv.FormatFlag)
	_ = p.handle.ObserveProperty(0, "time-pos/full", mpv.FormatDouble)
	_ = p.handle.ObserveProperty(0, "demuxer-cache-time", mpv.FormatDouble)
	err := p.handle.Initialize()

	if err != nil {
		fmt.Printf("[p.handle.Initialize]: %v\n", err)
		panic(err)

	}

	err = p.handle.Command([]string{"loadfile", p.URL})
	if err != nil {
		fmt.Printf("[p.handle.LoadFile]: %v (File: %v)\n", err, "balls")
		fmt.Printf("That didnt work.")
	}

	// Start the MPV event loop
	p.loop = true
	p.mutex.Unlock()

	p.startPlayerEventLoop()
	p.wg.Done()
}

func (p *Player) startPlayerEventLoop() {
	// Make sure to not change the underlying thread context
	r.LockOSThread()

	// Once the event loop is done we can release the thread back to go
	defer r.UnlockOSThread()

	for p.loop {
		// If the mpv handle was destroyed do not try and run the event loop
		if p.handle == nil {
			return
		}
		// MPV requires timeout to wait for events, wait for approx. 20ms between each event for accurate playback position
		ev := p.handle.WaitEvent(20)
		switch ev.EventID {
		// Property change events
		case mpv.EventPropertyChange:
			prop := ev.Property()
			value := prop.Data

			// Get pause property change
			if prop.Name == "pause" {
				p.paused = value.(int) != 0
			}

			// Get time position property change
			if prop.Name == "time-pos/full" {
				runtime.EventsEmit(p.ctx, PLAYER_POSITION, value)
			}

			if prop.Name == "demuxer-cache-time" && value != nil {
				runtime.EventsEmit(p.ctx, PLAYER_BUFFER, value.(float64))
			}
			// Runs when the file is initally loaded and ready to play
		case mpv.EventFileLoaded:
			fmt.Printf("File loaded and started playback\n")
			duration_val, err := p.handle.GetProperty("duration", mpv.FormatInt64)

			if err != nil {
				runtime.EventsEmit(p.ctx, PLAYER_GENERIC_ERROR, ev.Error)
			}
			runtime.EventsEmit(p.ctx, PLAYER_LOADED, duration_val.(int64))
			// Player was requested to seek, either internally or by the controller
		case mpv.EventSeek:
			println("Seek event requested")
			runtime.EventsEmit(p.ctx, PLAYER_LOADING)
			// Playback was restarted, either this happens when the player looses network then regains it, or after a seek was completed
		case mpv.EventPlaybackRestart:
			println("Playback restarted")
			runtime.EventsEmit(p.ctx, PLAYER_PLAYBACK_RESTART)
			// Regular log messages
		case mpv.EventLogMsg:
			msg := ev.LogMessage()
			fmt.Printf("player msg: %v\n", msg.Text)
			runtime.EventsEmit(p.ctx, PLAYER_MESSAGE, msg.Text)

			// Runs when the file is finished
		case mpv.EventEnd:
			ef := ev.EndFile()
			fmt.Println("end:", ef.EntryID, ef.Reason)
			if ef.Reason == mpv.EndFileEOF {
				p.loop = false
				return
			} else if ef.Reason == mpv.EndFileError {
				fmt.Println("error:", ef.Error)
			}
			// Runs when MPV player is shutdown
		case mpv.EventShutdown:
			fmt.Println("shutdown:", ev.EventID)
			p.loop = false
			return
		default:
			fmt.Println("event:", ev.EventID)
		}

		if ev.Error != nil {
			fmt.Println("error:", ev.Error)
			runtime.EventsEmit(p.ctx, PLAYER_GENERIC_ERROR, ev.Error)
		}
	}
}
