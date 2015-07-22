package commportal

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/FogCreek/victor"
	"github.com/groob/CommPhone/commportal"
	"github.com/groob/radigast/plugins"
)

type Phone struct {
	Username string
	Password string
	URL      string
}

func (c Phone) Register() []victor.HandlerDocPair {
	return []victor.HandlerDocPair{
		&victor.HandlerDoc{
			CmdHandler:     c.phoneLine,
			CmdName:        "phone",
			CmdDescription: "Prints the phone line or extension for a user",
			CmdUsage: []string{
				"[options] NAME",
				"-line _return direct line instead of extension_",
			},
		},
	}
}

func (c Phone) phoneLine(s victor.State) {
	var line bool
	phoneFlagSet := flag.NewFlagSet("phone", flag.ExitOnError)
	phoneFlagSet.BoolVar(&line, "line", false, "")
	phoneFlagSet.Parse(s.Fields())
	args := phoneFlagSet.Args()
	name := strings.Join(args, " ")
	if len(args) == 0 {
		msg := "You must add a name after the `phone` command."
		s.Chat().Send(s.Message().Channel().ID(), msg)
		return
	}
	c.phoneLineReply(s, name, line)

}

func (c Phone) phones() (commportal.Subscribers, error) {
	portal, err := commportal.New(c.Username, c.Password, c.URL)
	if err != nil {
		return nil, err
	}
	phones, err := portal.Phones()
	if err != nil {
		return nil, err
	}
	return phones, nil
}

func (c Phone) phoneLineReply(s victor.State, name string, line bool) {
	subscribers, err := c.phones()
	if err != nil {
		s.Chat().Send(s.Message().Channel().ID(), fmt.Sprintf("ERROR: `%s`", err))
		log.Println(err)
	}

	phones := subscribers.Find(name)
	for _, phone := range phones {
		var msg string
		if line {
			msg = fmt.Sprintf("%v %v", phone.Name.Name, phone.DirectoryNumber.Line)
		} else {
			msg = fmt.Sprintf("%v %v", phone.Name.Name, phone.IntercomDialingCode.Extension)
		}
		s.Chat().Send(s.Message().Channel().ID(), fmt.Sprintf("%v", msg))
	}
}

func init() {
	plugins.Add("phone", func() plugins.Registrator {
		return &Phone{}
	})
}
