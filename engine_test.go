package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWindow(t *testing.T) {
	err := Init()
	require.NoError(t, err)

	window, err := NewWindow()
	require.NoError(t, err)

	for window.ProcessEvents() {
		window.Blit()
	}

	window.Close()
}
