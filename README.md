# Radigast - Go Slack bot with configurable plugins.

Radigast is based on [https://github.com/FogCreek/victor](https://github.com/FogCreek/victor) which is a fork of [https://github.com/brettbuddin/victor](https://github.com/brettbuddin/victor)
and uses the Slack Real Time Messaging API.

# Configuration
Radigast loads configuration from a TOML file. 
The main configuration must be under `[radigast]` and the configuration for each plugin must be under `[plugin_name]`

See [config.toml.sample](https://github.com/groob/radigast/blob/master/config.toml.sample) for a complete example.

# How to use it
Run `radigast -config radigast.toml` to connect to slack.

# Plugins
Radigast supports two types of plugins:
The first type of plugin is a Go plugin that must be compiled with the Radigast cli tool. See more under `Developing Go Plugins`.
Note: The plugins in /plugins/all are not currently compiled into the released binary as not everyone will find them useful.

The second type of plugin uses JSON-RPC and can be written in most languages.
There is a python example under [examples-rpc](https://github.com/groob/radigast/tree/master/example-rpc/python)

An RPC plugin must provide the following methods: 

```Go
 Name() string
 Description() string
 Usage() []string
 Handle(args) string
```
 The args in Handle(args) will be of the following struct type:

```Go
 type Args struct {
 	// Chat user calling the plugin.
 	User string
 	// The arguments a user passes to the bot.
 	Fields []string
 }

```

RPC plugins must be installed in a directory specified in the config file. Each plugin name must be added to the `rpcplugins` array in the config.


## Developing Go Plugins
Radigast uses the same plugin model as [telegraf](https://github.com/influxdb/telegraf) 

* A plugin must conform to the `plugins.Registrator` interface
* Plugins should call plugins.Add in their init function to register themselves
* To be available to the Radigast command, plugins must be added to the `github.com/groob/radigast/plugins/all/all.go` file.
* A plugin will only be configured by radigast if there is a `[plugin_name]` section in the config file

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
func (h Hello) Register() []victor.HandlerDocPair {
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
	s.Reply(msg)
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
