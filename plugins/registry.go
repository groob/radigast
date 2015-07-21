package plugins

import (
	"github.com/FogCreek/victor"
)

// Registrator is the interface that a plugin must implement.
type Registrator interface {
	Register() victor.HandlerDocPair
}

// Plugin returns a plugin registrator
type Plugin func() Registrator

// map of registered plugins
var Plugins = map[string]Plugin{}

// Add takes a command name and Plugin function
// Add should be used in the plugins init() function
func Add(command string, plugin Plugin) {
	Plugins[command] = plugin

}
