package model

type (
	SeshSessions struct {
		// catalog of the sessions
		Directory SeshSessionMap
		// unique identifiers of the sessions ordered
		OrderedIndex []string
	}

	SeshSessionMap map[string]SeshSession
	SeshWindowMap  map[string]WindowConfig

	SeshSession struct {
		Src  string // The source of the session (config, tmux, zoxide, tmuxinator)
		Name string // The display name
		Path string // The absolute directory path

		StartupCommand        string         // The command to run when the session is started
		PreviewCommand        string         // The command to run when the session is previewed
		DisableStartupCommand bool           // Ignore the default startup command if present
		SkipDefaultWindow     bool           // Remove tmux's default window when creating configured windows
		SourcePath            string         // Source config file path
		Tmuxinator            string         // Name of the tmuxinator config
		Attached              int            // Whether the session is currently attached
		Windows               int            // The number of windows in the session
		WindowConfigs         []WindowConfig // The windows used in session config
		WindowNames           []string       // The names of the windows in session config
		PaneNames             []string       // The names of the panes used in window config
		Score                 float64        // The score of the session (from Zoxide)
	}

	SeshSrcs struct {
		Config     bool
		Tmux       bool
		Tmuxinator bool
		Zoxide     bool
	}
)
