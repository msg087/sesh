package lister

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func ConfigKey(name string) string {
	return fmt.Sprintf("config:%s", name)
}

func listConfig(l *RealLister) (model.SeshSessions, error) {
	orderedIndex := make([]string, 0)
	directory := make(model.SeshSessionMap)
	for _, session := range l.config.SessionConfigs {
		if session.Name != "" {
			key := ConfigKey(session.Name)
			orderedIndex = append(orderedIndex, key)
			path, err := l.home.ExpandPath(session.Path)
			if err != nil {
				return model.SeshSessions{}, fmt.Errorf("couldn't expand home: %q", err)
			}

			if session.StartupCommand != "" && session.DisableStartCommand {
				return model.SeshSessions{}, fmt.Errorf("startup_command and disable_start_command are mutually exclusive")
			}

			windowConfigs := make([]model.WindowConfig, 0, len(session.Windows))
			for _, windowName := range session.Windows {
				for _, windowConfig := range l.config.WindowConfigs {
					if windowConfig.Name == windowName {
						windowConfigs = append(windowConfigs, windowConfig)
						break
					}
				}
			}

			directory[key] = model.SeshSession{
				Src:                   "config",
				Name:                  session.Name,
				Path:                  path,
				SourcePath:            session.SourcePath,
				StartupCommand:        session.StartupCommand,
				PreviewCommand:        session.PreviewCommand,
				DisableStartupCommand: session.DisableStartCommand,
				SkipDefaultWindow:     session.SkipDefaultWindow,
				Tmuxinator:            session.Tmuxinator,
				WindowConfigs:         windowConfigs,
				WindowNames:           session.Windows,
				PaneNames:             session.Panes,
			}
		}
	}
	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindConfigSession(name string) (model.SeshSession, bool) {
	key := ConfigKey(name)
	sessions, _ := listConfig(l)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	} else {
		return model.SeshSession{}, false
	}
}
