package tmux

func (t *RealTmux) KillWindow(targetWindow string) (string, error) {
	return t.shell.Cmd("tmux", "kill-window", "-t", targetWindow)
}
