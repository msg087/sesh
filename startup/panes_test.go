package startup

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/replacer"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
)

func TestExecBuildsPaneLayout(t *testing.T) {
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockHome := new(home.MockHome)
	mockReplacer := new(replacer.MockReplacer)

	s := &RealStartup{
		lister: mockLister,
		tmux:   mockTmux,
		config: model.Config{
			WindowConfigs: []model.WindowConfig{{
				Name:  "api",
				Path:  "/repo/api",
				Panes: []string{"editor", "terminal", "side"},
			}},
			PaneConfigs: []model.PaneConfig{
				{Name: "editor", StartupScript: "nvim ."},
				{Name: "terminal", Split: "v", Size: 12, SizeMode: "lines", StartupScript: "clear"},
				{Name: "side", Split: "h", Size: 30, SizeMode: "percent", StartupScript: "htop", Panes: []string{"child"}},
				{Name: "child", Split: "v", Size: 20, SizeMode: "percent", StartupScript: "echo child"},
			},
		},
		home:     mockHome,
		replacer: mockReplacer,
	}

	session := model.SeshSession{
		Name:        "JSR",
		Path:        "/repo/jsr",
		WindowNames: []string{"api"},
	}

	mockLister.On("FindConfigSession", "JSR").Return(model.SeshSession{}, false)
	mockLister.On("FindConfigWildcard", "/repo/jsr").Return(model.WildcardConfig{}, false)
	mockHome.On("ExpandPath", "/repo/api").Return("/repo/api", nil)
	mockTmux.On("NewWindow", "/repo/api", "api").Return("", nil)
	mockTmux.On("SplitWindow", "JSR:api.0", "/repo/api", "v", 12, "lines").Return("%1", nil)
	mockTmux.On("SendKeys", "%1", "clear").Return("", nil)
	mockTmux.On("SplitWindow", "%1", "/repo/api", "h", 30, "percent").Return("%2", nil)
	mockTmux.On("SendKeys", "%2", "htop").Return("", nil)
	mockTmux.On("SplitWindow", "%2", "/repo/api", "v", 20, "percent").Return("%3", nil)
	mockTmux.On("SendKeys", "%3", "echo child").Return("", nil)
	mockTmux.On("SendKeys", "JSR:api.0", "cd \"/repo/api\" && nvim .").Return("", nil)
	mockTmux.On("NextWindow").Return("", nil)
	mockTmux.On("SelectPane", 0, 0).Return("", nil)

	command, err := s.Exec(session)
	assert.NoError(t, err)
	assert.Empty(t, command)
}
