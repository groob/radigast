package hello

import (
	"fmt"

	"github.com/FogCreek/victor"
	"github.com/groob/radigast/plugins"
)

// Configuration struct
// toml will unmarshal any options provided under [hello] in
// radigast.toml
type Hello struct {
	// AnOption      string
	// AnotherOption string
}

// Register implements plugins.Registrator
func (h Hello) Register() victor.HandlerDocPair {
	return &victor.HandlerDoc{
		CmdHandler:     h.helloFunc,
		CmdName:        "hello",
		CmdDescription: "reply back with the user name",
		CmdUsage:       []string{"NAME"},
	}
}

// Bot Handler
// write your plugin logic here.
func (h Hello) helloFunc(s victor.State) {
	msg := fmt.Sprintf("Hello %s!", s.Message().User().Name())
	s.Chat().Send(s.Message().Channel().ID(), msg)
}

func init() {
	// register the plugin
	plugins.Add("hello", func() plugins.Registrator {
		return &Hello{}
	})
}
