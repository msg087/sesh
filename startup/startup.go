package startup

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
)

type Startup interface {
	Exec(session model.SeshSession) (string, error)
}

type RealStartup struct {
	lister   lister.Lister
	tmux     tmux.Tmux
	config   model.Config
	home     home.Home
	replacer replacer.Replacer
}

func NewStartup(
	config model.Config, lister lister.Lister, tmux tmux.Tmux, home home.Home, replacer replacer.Replacer,
) Startup {
	return &RealStartup{lister, tmux, config, home, replacer}
}

func (s *RealStartup) Exec(session model.SeshSession) (string, error) {
	strategies := []func(*RealStartup, model.SeshSession) (string, error){
		configStrategy,
		configWildcardStartupStrategy,
		defaultConfigStrategy,
	}

	windows := make(map[string]model.SeshWindowMap)
	for _, window := range s.config.WindowConfigs {
		var path string = ""
		var err error = nil
		if window.Path != "" {
			if strings.HasPrefix(window.Path, "~") {
				path, err = s.home.ExpandPath(window.Path)
				if err != nil {
					return "", fmt.Errorf("couldn't expand home: %q", err)
				}
			} else {
				path = window.Path
			}
		}

		if windows[window.SourcePath] == nil {
			windows[window.SourcePath] = make(model.SeshWindowMap)
		}
		windows[window.SourcePath][window.Name] = model.WindowConfig{
			Name:          window.Name,
			Path:          path,
			StartupScript: window.StartupScript,
			Panes:         window.Panes,
			SourcePath:    window.SourcePath,
		}
	}
	paneConfigs := make(map[string]map[string]model.PaneConfig)
	for _, pane := range s.config.PaneConfigs {
		if paneConfigs[pane.SourcePath] == nil {
			paneConfigs[pane.SourcePath] = make(map[string]model.PaneConfig)
		}
		paneConfigs[pane.SourcePath][pane.Name] = pane
	}
	skipDefaultWindow := session.SkipDefaultWindow && len(session.WindowNames) > 0

	for i, window := range session.WindowNames {
		windowConfig, ok := resolveWindowConfig(session.SourcePath, window, windows)
		if !ok {
			return "", fmt.Errorf("window %s is not defined in config", window)
		}
		if windowConfig.Path == "" {
			path, err := s.home.ExpandPath(session.Path)
			if err != nil {
				return "", fmt.Errorf("couldn't expand home: %q", err)
			}
			windowConfig.Path = path
		} else if strings.HasPrefix(windowConfig.Path, "~") {
			path, err := s.home.ExpandPath(windowConfig.Path)
			if err != nil {
				return "", fmt.Errorf("couldn't expand home: %q", err)
			}
			windowConfig.Path = path
		} else if !filepath.IsAbs(windowConfig.Path) {
			windowConfig.Path = filepath.Join(session.Path, windowConfig.Path)
		}

		if !(session.SkipDefaultWindow && i == 0) {
			if ret, err := s.tmux.NewWindow(windowConfig.Path, windowConfig.Name); err != nil {
				return ret, err
			}
		}

		paneNames := windowConfig.Panes
		if len(paneNames) == 0 {
			paneNames = session.PaneNames
		}
		if len(paneNames) == 0 {
			targetPane := fmt.Sprintf("%s:%s.0", session.Name, windowConfig.Name)
			if windowConfig.StartupScript != "" {
				if ret, err := s.tmux.SendKeys(targetPane, windowConfig.StartupScript); err != nil {
					return ret, err
				}
			}
			if _, err := s.tmux.SelectPaneTarget(targetPane); err != nil {
				return "", err
			}
			continue
		}

		if ret, err := s.buildWindowPanes(session, windowConfig, paneNames, paneConfigs); err != nil {
			return ret, err
		}

		if _, err := s.tmux.SelectPaneTarget(fmt.Sprintf("%s:%s.0", session.Name, windowConfig.Name)); err != nil {
			return "", err
		}
	}

	if !skipDefaultWindow {
		s.tmux.NextWindow()
	}
	for _, strategy := range strategies {
		if command, err := strategy(s, session); err != nil {
			return "", fmt.Errorf("failed to determine startup command: %w", err)
		} else if command != "" {
			s.tmux.SendKeys(session.Name, command)
			return fmt.Sprintf("executing startup command: %s", command), nil
		}
	}

	return "", nil // no command to run
}
