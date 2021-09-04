package engine

import (
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func Init() error {
	runtime.LockOSThread()

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_EVENTS)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %w", err)
	}

	err = sdl.Init(mix.INIT_MP3)
	if err != nil {
		return fmt.Errorf("could not initialize MP3 playback: %w", err)
	}

	return nil
}

func Close() {
	mix.Quit()
	sdl.Quit()
}
