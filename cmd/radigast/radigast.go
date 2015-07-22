package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/FogCreek/victor"
	"github.com/FogCreek/victor/pkg/chat/slackRealtime"
	"github.com/FogCreek/victor/pkg/events"
	"github.com/groob/radigast"
	_ "github.com/groob/radigast/plugins/all"
)

var fConfig = flag.String("config", "", "configuration file to load")
var fVersion = flag.Bool("version", false, "display the version")

// Version string
var Version = "unreleased"

func main() {
	flag.Parse()

	if *fVersion {
		fmt.Printf("Radigast Slack Bot - Version %s\n", Version)
		return
	}
	var (
		config *radigast.Config
		err    error
	)
	if *fConfig != "" {
		config, err = radigast.LoadConfig(*fConfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No config file specified.")
	}

	bot := victor.New(victor.Config{
		ChatAdapter:   "slackRealtime",
		AdapterConfig: slackRealtime.NewConfig(config.SlackToken),
		Name:          config.BotName,
	})

	// load radigast plugins
	config.LoadPlugins(bot)
	config.LoadRPCPlugins(bot)

	// Enable in-chat help command
	bot.EnableHelpCommand()

	//  run the bot
	runBot(bot)
}

// run and recover from panic during bot.Stop()
// the upstream slack adapter has not implemented bot.Stop() yet
func runBot(bot victor.Robot) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("bot.Stop() exited with panic: ", e)
			os.Exit(0)
		}
	}()

	bot.Run()
	go monitorErrors(bot.ChatErrors())
	go monitorEvents(bot.ChatEvents())
	// keep the process (and bot) alive
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	bot.Stop()
}

func monitorErrors(errorChannel <-chan events.ErrorEvent) {
	for {
		err, ok := <-errorChannel
		if !ok {
			return
		}
		if err.IsFatal() {
			log.Panic(err.Error())
		}
		log.Println("Chat Adapter Error Event:", err.Error())
	}
}

func monitorEvents(eventsChannel chan events.ChatEvent) {
	for {
		e, ok := <-eventsChannel
		if !ok {
			return
		}
		log.Printf("Chat Event: %+v", e)
	}
}
