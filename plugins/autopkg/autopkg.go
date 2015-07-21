package autopkg

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/FogCreek/victor"
	"github.com/groob/radigast/plugins"
)

type Autopkg struct {
	Path         string
	AllowedUsers []string
}

func (a Autopkg) Register() victor.HandlerDocPair {
	// Allow everyone or just a specific group of users?
	var handler victor.HandlerFunc
	if len(a.AllowedUsers) == 0 {
		handler = a.autopkgFunc
	} else {
		handler = victor.OnlyAllow(a.AllowedUsers, a.autopkgFunc)
	}

	return &victor.HandlerDoc{
		CmdHandler:     handler,
		CmdName:        "autopkg",
		CmdDescription: "check for new versions of software and add to munki",
		CmdUsage:       []string{""},
	}
}

func (a Autopkg) autopkgFunc(s victor.State) {
	autopkgCmd := exec.Command(a.Path)
	for _, arg := range s.Fields() {
		autopkgCmd.Args = append(autopkgCmd.Args, arg)
	}
	var out bytes.Buffer
	autopkgCmd.Stdout = &out
	err := autopkgCmd.Run()
	if err != nil {
		log.Println(err)
	}
	s.Chat().Send(s.Message().Channel().ID(), out.String())
}

func init() {
	plugins.Add("autopkg", func() plugins.Registrator {
		return &Autopkg{}
	})
}
