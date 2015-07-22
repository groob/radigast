package radigast

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/FogCreek/victor"
	"github.com/groob/radigast/plugins"
	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

// Config holds the radigast configuration
type Config struct {
	SlackToken string
	BotName    string

	plugins map[string]*ast.Table
}

// Invalid toml config
var ErrInvalidConfig = errors.New("invalid configuration")

// LoadConfig unmarshalls toml config file for radigast and all plugins
func LoadConfig(path string) (*Config, error) {
	c := &Config{
		plugins: make(map[string]*ast.Table),
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tbl, err := toml.Parse(data)
	if err != nil {
		return nil, err
	}

	for name, val := range tbl.Fields {
		subtbl, ok := val.(*ast.Table)
		if !ok {
			return nil, ErrInvalidConfig
		}

		switch name {
		case "slackbot":
			err := toml.UnmarshalTable(subtbl, c)
			if err != nil {
				return nil, err
			}

		default:
			c.plugins[name] = subtbl
		}
	}

	return c, nil
}

// LoadPlugins registers victor.Handlers with radigast
// and adds any additional configuration to the plugin
func (c Config) LoadPlugins(r victor.Robot) {
	log.Printf("Loading %v plugins\n", len(c.plugins))

	for name := range plugins.Plugins {
		registrator, ok := plugins.Plugins[name]
		if !ok {
			log.Printf("Undefined but requested plugin: %s", name)
		}

		plugin := registrator()

		// apply configuration to the plugin
		c.applyPluginConfig(name, plugin, r)

	}

	// Add default handler to show "unrecognized command" on "command" messages
	r.SetDefaultHandler(defaultFunc)
}

func (c Config) applyPluginConfig(name string, plugin plugins.Registrator, r victor.Robot) {
	if tbl, ok := c.plugins[name]; ok {
		err := toml.UnmarshalTable(tbl, plugin)
		if err != nil {
			log.Printf("Couldn't Unmarshal config for %v\n", name)
		}

		// register a plugin's handlers with bot
		handlers := plugin.Register()
		for _, handler := range handlers {
			r.HandleCommand(handler)
		}

		log.Printf("Loaded %s\n", name)
	}
}

func defaultFunc(s victor.State) {
	s.Chat().Send(s.Message().Channel().ID(),
		"Unrecognized command. Type `help` to see supported commands.")
}
