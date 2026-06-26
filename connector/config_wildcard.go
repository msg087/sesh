package connector

import (
	"github.com/joshmedeski/sesh/v2/model"
)

func configWildcardStrategy(c *RealConnector, name string) (model.Connection, error) {
	wc, found := c.lister.FindConfigWildcard(name)
	if !found {
		return model.Connection{Found: false}, nil
	}

	path, err := c.home.ExpandPath(name)
	if err != nil {
		return model.Connection{}, err
	}

	isDir, absPath := c.dir.Dir(path)
	if !isDir {
		return model.Connection{Found: false}, nil
	}

	nameFromPath, err := c.namer.Name(absPath)
	if err != nil {
		return model.Connection{}, err
	}

	panes := wc.Panes
	if len(panes) == 0 {
		panes = c.config.DefaultSessionConfig.Panes
	}

	return model.Connection{
		Found:       true,
		New:         true,
		AddToZoxide: true,
		Session: model.SeshSession{
			Src:                   "config_wildcard",
			Name:                  nameFromPath,
			Path:                  absPath,
			SourcePath:            wc.SourcePath,
			WindowNames:           wc.Windows,
			PaneNames:             panes,
			DisableStartupCommand: wc.DisableStartCommand,
			SkipDefaultWindow:     wc.SkipDefaultWindow,
		},
	}, nil
}
