# Radigast - Go Slack bot with configurable plugins.

Radigast is based on [https://github.com/FogCreek/victor](https://github.com/FogCreek/victor) which is a fork of [https://github.com/brettbuddin/victor](https://github.com/brettbuddin/victor)
and uses the Slack Real Time Messaging API.

# Configuration
Radigast loads configuration from a TOML file. 
The main configuration must be under `[radigast]` and the configuration for each plugin must be under `[plugin_name]`

See [config.toml.sample]() for a complete example.

# How to use it
Run `radigast -config radigast.toml` to connect to slack.


# Plugins

## Developing Plugins
Radigast uses the same plugin model as [telegraf](https://github.com/influxdb/telegraf) 

* A plugin must conform to the `plugins.Registrator` interface
* Plugins should call plugins.Add in their init function to register themselves
* To be available to the Radigast command, plugins must be added to the `github.com/groob/radigast/plugins/all/all.go` file.

## Plugin Interface
```Go
type Registrator interface {
	Register() []victor.HandlerDocPair
}
```
## Plugin example

```Go
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
	return []victor.HandlerDocPair{
		&victor.HandlerDoc{
			CmdHandler:     h.helloFunc,
			CmdName:        "hello",
			CmdDescription: "reply back with the user name",
			CmdUsage:       []string{"NAME"},
		},
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
```

# Example usage
![](http://i.imgur.com/S9zF8Jc.png)
