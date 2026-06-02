package configurator

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

func TestPaneConfigParsing(t *testing.T) {
	input := []byte(`
[[pane]]
name = "editor"
startup_script = "nvim ."

[[pane]]
name = "terminal"
split = "v"
size = 12
size_mode = "lines"
startup_script = "clear"

[[pane]]
name = "mini_split_term"
split = "v"
size = 12
size_mode = "lines"
panes = ["split_terminal"]

[[pane]]
name = "split_terminal"
split = "h"
size = 20
size_mode = "percent"

[default_session]
panes = ["terminal"]

[[window]]
name = "api"
panes = ["editor", "terminal"]
`)

	config := model.Config{}
	err := toml.Unmarshal(input, &config)
	assert.NoError(t, err)
	assert.Len(t, config.PaneConfigs, 4)
	assert.Equal(t, "editor", config.PaneConfigs[0].Name)
	assert.Equal(t, "nvim .", config.PaneConfigs[0].StartupScript)
	assert.Equal(t, "terminal", config.PaneConfigs[1].Name)
	assert.Equal(t, "v", config.PaneConfigs[1].Split)
	assert.Equal(t, 12, config.PaneConfigs[1].Size)
	assert.Equal(t, "lines", config.PaneConfigs[1].SizeMode)
	assert.Equal(t, []string{"split_terminal"}, config.PaneConfigs[2].Panes)
	assert.Equal(t, "split_terminal", config.PaneConfigs[3].Name)
	assert.Equal(t, "h", config.PaneConfigs[3].Split)
	assert.Equal(t, 20, config.PaneConfigs[3].Size)
	assert.Equal(t, "percent", config.PaneConfigs[3].SizeMode)
	assert.Equal(t, []string{"terminal"}, config.DefaultSessionConfig.Panes)
	assert.Equal(t, []string{"editor", "terminal"}, config.WindowConfigs[0].Panes)
}

func TestPaneImportMerge(t *testing.T) {
	mainConfig := model.Config{}
	importConfig := model.Config{
		PaneConfigs: []model.PaneConfig{{Name: "editor"}},
	}
	mainConfig.PaneConfigs = append(mainConfig.PaneConfigs, importConfig.PaneConfigs...)
	assert.Len(t, mainConfig.PaneConfigs, 1)
	assert.Equal(t, "editor", mainConfig.PaneConfigs[0].Name)
}
