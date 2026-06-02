package tmux

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/oswrap"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/stretchr/testify/assert"
)

func TestClientSize(t *testing.T) {
	t.Run("uses tmux client size when attached", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockShell := new(shell.MockShell)
		tmux := NewTmux(mockOs, mockShell, "")

		mockOs.On("Getenv", "TMUX").Return("/tmp/tmux-1/default,1,0")
		mockShell.On("Cmd", "tmux", "display-message", "-p", "#{client_width}").Return("120", nil)
		mockShell.On("Cmd", "tmux", "display-message", "-p", "#{client_height}").Return("40", nil)

		cols, lines, err := tmux.ClientSize()
		assert.NoError(t, err)
		assert.Equal(t, 120, cols)
		assert.Equal(t, 40, lines)
	})

	t.Run("falls back outside tmux", func(t *testing.T) {
		mockOs := new(oswrap.MockOs)
		mockShell := new(shell.MockShell)
		tmux := NewTmux(mockOs, mockShell, "")

		mockOs.On("Getenv", "TMUX").Return("")

		orig := terminalSize
		defer func() { terminalSize = orig }()
		terminalSize = func() (int, int) { return 120, 40 }

		cols, lines, err := tmux.ClientSize()
		assert.NoError(t, err)
		assert.Equal(t, 120, cols)
		assert.Equal(t, 40, lines)
	})
}
