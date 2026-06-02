package startup

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func (s *RealStartup) buildWindowPanes(session model.SeshSession, windowConfig model.WindowConfig, paneNames []string, paneConfigs map[string]map[string]model.PaneConfig) (string, error) {
	if len(paneNames) == 0 {
		return "", nil
	}

	basePane := fmt.Sprintf("%s:%s.0", session.Name, windowConfig.Name)
	firstPane, err := resolvePaneConfig(session.SourcePath, paneNames[0], paneConfigs)
	if err != nil {
		return "", err
	}

	currentTarget := basePane
	if ret, err := s.applyPaneConfig(session, windowConfig, basePane, firstPane, paneConfigs); err != nil {
		return ret, err
	} else {
		currentTarget = ret
	}

	for _, paneName := range paneNames[1:] {
		paneConfig, err := resolvePaneConfig(session.SourcePath, paneName, paneConfigs)
		if err != nil {
			return "", err
		}
		if ret, err := s.buildPaneFromTarget(session, windowConfig, currentTarget, paneConfig, paneConfigs); err != nil {
			return ret, err
		} else {
			currentTarget = ret
		}
	}

	return "", nil
}

func resolvePaneConfig(sourcePath, name string, paneConfigs map[string]map[string]model.PaneConfig) (model.PaneConfig, error) {
	if configs := paneConfigs[sourcePath]; configs != nil {
		if paneConfig, ok := configs[name]; ok {
			return paneConfig, nil
		}
	}

	var (
		found    model.PaneConfig
		foundOne bool
	)
	for _, configs := range paneConfigs {
		if paneConfig, ok := configs[name]; ok {
			if foundOne {
				return model.PaneConfig{}, fmt.Errorf("pane %s is ambiguous across imported config files", name)
			}
			found = paneConfig
			foundOne = true
		}
	}
	if !foundOne {
		return model.PaneConfig{}, fmt.Errorf("pane %s is not defined in config", name)
	}
	return found, nil
}

func resolveWindowConfig(sourcePath, name string, windows map[string]model.SeshWindowMap) (model.WindowConfig, bool) {
	if sourceWindows := windows[sourcePath]; sourceWindows != nil {
		if windowConfig, ok := sourceWindows[name]; ok {
			return windowConfig, true
		}
	}

	var (
		found    model.WindowConfig
		foundOne bool
	)
	for _, sourceWindows := range windows {
		if windowConfig, ok := sourceWindows[name]; ok {
			if foundOne {
				return model.WindowConfig{}, false
			}
			found = windowConfig
			foundOne = true
		}
	}
	return found, foundOne
}

func (s *RealStartup) buildPaneFromTarget(session model.SeshSession, windowConfig model.WindowConfig, targetPane string, paneConfig model.PaneConfig, paneConfigs map[string]map[string]model.PaneConfig) (string, error) {
	paneID, err := s.tmux.SplitWindow(targetPane, windowConfig.Path, paneConfig.Split, paneConfig.Size, paneConfig.SizeMode)
	if err != nil {
		return paneID, err
	}
	return s.applyPaneConfig(session, windowConfig, paneID, paneConfig, paneConfigs)
}

func (s *RealStartup) applyPaneConfig(session model.SeshSession, windowConfig model.WindowConfig, targetPane string, paneConfig model.PaneConfig, paneConfigs map[string]map[string]model.PaneConfig) (string, error) {
	currentTarget := targetPane
	for _, childName := range paneConfig.Panes {
		childConfig, err := resolvePaneConfig(session.SourcePath, childName, paneConfigs)
		if err != nil {
			return "", err
		}
		if ret, err := s.buildPaneFromTarget(session, windowConfig, currentTarget, childConfig, paneConfigs); err != nil {
			return ret, err
		} else {
			currentTarget = ret
		}
	}

	command := paneConfig.StartupScript
	if command == "" && targetPane == fmt.Sprintf("%s:%s.0", session.Name, windowConfig.Name) {
		command = windowConfig.StartupScript
	}
	if command != "" && targetPane == fmt.Sprintf("%s:%s.0", session.Name, windowConfig.Name) && windowConfig.Path != "" {
		command = fmt.Sprintf("cd %q && %s", windowConfig.Path, command)
	}
	if command != "" {
		if ret, err := s.tmux.SendKeys(targetPane, command); err != nil {
			return ret, err
		}
	}

	return currentTarget, nil
}
