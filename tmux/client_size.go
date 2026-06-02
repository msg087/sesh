package tmux

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

var terminalSize = func() (int, int) {
	cols, lines, err := term.GetSize(os.Stdout.Fd())
	if err == nil {
		return cols, lines
	}

	cols, lines, err = term.GetSize(os.Stdin.Fd())
	if err == nil {
		return cols, lines
	}

	cols, lines, err = term.GetSize(os.Stderr.Fd())
	if err == nil {
		return cols, lines
	}

	return 80, 24
}

func (t *RealTmux) ClientSize() (int, int, error) {
	if t.os.Getenv("TMUX") != "" {
		width, err := t.shell.Cmd(t.bin, "display-message", "-p", "#{client_width}")
		if err == nil {
			height, err := t.shell.Cmd(t.bin, "display-message", "-p", "#{client_height}")
			if err == nil {
				cols, err := strconv.Atoi(strings.TrimSpace(width))
				if err != nil {
					return 0, 0, err
				}
				lines, err := strconv.Atoi(strings.TrimSpace(height))
				if err != nil {
					return 0, 0, err
				}
				return cols, lines, nil
			}
		}
	}

	cols, lines := terminalSize()
	return cols, lines, nil
}
