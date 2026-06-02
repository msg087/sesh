package tmux

import "strconv"

func (t *RealTmux) SplitWindow(target string, startDir string, split string, size int, sizeMode string) (string, error) {
	args := []string{"split-window", "-P", "-F", "#{pane_id}"}
	if split == "h" {
		args = append(args, "-h")
	} else if split == "v" {
		args = append(args, "-v")
	}
	if size > 0 {
		if sizeMode == "percent" {
			args = append(args, "-p", strconv.Itoa(size))
		} else {
			args = append(args, "-l", strconv.Itoa(size))
		}
	}
	if target != "" {
		args = append(args, "-t", target)
	}
	if startDir != "" {
		args = append(args, "-c", startDir)
	}
	return t.shell.Cmd(t.bin, args...)
}
