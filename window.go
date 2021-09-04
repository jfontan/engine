package engine

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	window  *sdl.Window
	context sdl.GLContext
}

func NewWindow() (*Window, error) {
	window, err := sdl.CreateWindow(
		"zooloo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_OPENGL,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create window: %w", err)
	}

	context, err := window.GLCreateContext()
	if err != nil {
		_ = window.Destroy()
		return nil, fmt.Errorf("could not create GL context: %w", err)
	}

	err = gl.Init()
	if err != nil {
		sdl.GLDeleteContext(context)
		_ = window.Destroy()
		return nil, fmt.Errorf("could not initialize OpenGL: %w", err)
	}

	gl.Viewport(0, 0, 800, 600)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	return &Window{
		window:  window,
		context: context,
	}, nil
}

func (w *Window) Close() error {
	sdl.GLDeleteContext(w.context)
	return w.window.Destroy()
}

func (w *Window) Blit() {
	w.window.GLSwap()
}

func (w *Window) ProcessEvents() bool {
	for {
		event := sdl.PollEvent()
		if event == nil {
			return true
		}

		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.Keycode(sdl.K_ESCAPE) {
				return false
			}
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_RESIZED {
				println("resized", t.Data1, t.Data2)
				gl.Viewport(0, 0, t.Data1, t.Data2)
			}
		}
	}
}
