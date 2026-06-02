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

func TestExecSkipsDefaultWindow(t *testing.T) {
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockHome := new(home.MockHome)
	mockReplacer := new(replacer.MockReplacer)

	s := &RealStartup{
		lister:   mockLister,
		tmux:     mockTmux,
		config:   model.Config{WindowConfigs: []model.WindowConfig{{Name: "api", StartupScript: "clear"}, {Name: "ns", StartupScript: "echo ns"}, {Name: "iac", StartupScript: "echo iac"}}},
		home:     mockHome,
		replacer: mockReplacer,
	}

	session := model.SeshSession{
		Name:              "JSR",
		Path:              "/repo/jsr-netsuite",
		WindowNames:       []string{"api", "ns", "iac"},
		SkipDefaultWindow: true,
	}

	mockLister.On("FindConfigSession", "JSR").Return(model.SeshSession{}, false)
	mockLister.On("FindConfigWildcard", "/repo/jsr-netsuite").Return(model.WildcardConfig{}, false)
	mockHome.On("ExpandPath", "/repo/jsr-netsuite").Return("/repo/jsr-netsuite", nil)
	mockTmux.On("NewWindow", "/repo/jsr-netsuite", "ns").Return("", nil)
	mockTmux.On("NewWindow", "/repo/jsr-netsuite", "iac").Return("", nil)
	mockTmux.On("SendKeys", "JSR", "clear").Return("", nil)
	mockTmux.On("SendKeys", "JSR", "echo ns").Return("", nil)
	mockTmux.On("SendKeys", "JSR", "echo iac").Return("", nil)
	mockTmux.On("SelectPane", 0, 0).Return("", nil)
	mockTmux.On("SelectWindow", "JSR:api").Return("", nil)

	command, err := s.Exec(session)
	assert.NoError(t, err)
	assert.Empty(t, command)
}

func TestExecResolvesRelativeWindowPathFromSessionPath(t *testing.T) {
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockHome := new(home.MockHome)
	mockReplacer := new(replacer.MockReplacer)

	s := &RealStartup{
		lister:   mockLister,
		tmux:     mockTmux,
		config:   model.Config{WindowConfigs: []model.WindowConfig{{Name: "api", Path: "./go_jsr_api"}}},
		home:     mockHome,
		replacer: mockReplacer,
	}

	session := model.SeshSession{
		Name:        "JSR",
		Path:        "/home/testuser/code_wsl/jsr-netsuite",
		WindowNames: []string{"api"},
	}

	mockLister.On("FindConfigSession", "JSR").Return(model.SeshSession{}, false)
	mockLister.On("FindConfigWildcard", "/home/testuser/code_wsl/jsr-netsuite").Return(model.WildcardConfig{}, false)
	mockTmux.On("NewWindow", "/home/testuser/code_wsl/jsr-netsuite/go_jsr_api", "api").Return("", nil)
	mockTmux.On("NextWindow").Return("", nil)
	mockTmux.On("SelectPane", 0, 0).Return("", nil)

	command, err := s.Exec(session)
	assert.NoError(t, err)
	assert.Empty(t, command)
}
