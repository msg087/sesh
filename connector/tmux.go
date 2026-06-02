package connector

import "github.com/joshmedeski/sesh/v2/model"

func tmuxStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindTmuxSession(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}
	return model.Connection{
		Found:       true,
		Session:     session,
		New:         false,
		AddToZoxide: true,
	}, nil
}

func connectToTmux(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	if connection.New {
		initialWindow := ""
		if connection.Session.SkipDefaultWindow && len(connection.Session.WindowNames) > 0 {
			initialWindow = connection.Session.WindowNames[0]
		}
		width, height, err := c.tmux.ClientSize()
		if err != nil {
			return "", err
		}
		if _, err := c.tmux.NewDetachedSession(connection.Session.Name, connection.Session.Path, initialWindow, width, height); err != nil {
			return "", err
		}
		if opts.Command != "" {
			if _, err := c.tmux.SendKeys(connection.Session.Name, opts.Command); err != nil {
				return "", err
			}
		} else {
			if _, err := c.startup.Exec(connection.Session); err != nil {
				return "", err
			}
		}
	}
	return c.tmux.SwitchOrAttach(connection.Session.Name, opts)
}
