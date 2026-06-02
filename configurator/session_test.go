package configurator

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

func TestSessionConfigParsing(t *testing.T) {
	input := []byte(`
[[session]]
name = "JSR"
path = "~/code_wsl/jsr-netsuite"
windows = ["api", "ns", "iac"]
skip_default_window = true
`)
	config := model.Config{}
	err := toml.Unmarshal(input, &config)
	assert.NoError(t, err)
	assert.Len(t, config.SessionConfigs, 1)
	assert.Equal(t, "JSR", config.SessionConfigs[0].Name)
	assert.True(t, config.SessionConfigs[0].SkipDefaultWindow)
}
