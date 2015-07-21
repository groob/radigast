package warranty

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/FogCreek/victor"
	"github.com/groob/radigast/plugins"
)

// Warranty plugin for radigast
type Warranty struct {
	Path string
}

// Register returns the victor handler
func (w Warranty) Register() victor.HandlerDocPair {
	return &victor.HandlerDoc{
		CmdHandler:     w.warrantyFunc,
		CmdName:        "warranty",
		CmdDescription: "Check Mac hardware warranty",
		CmdUsage:       []string{"SERIAL"},
	}
}

func (w Warranty) warrantyFunc(s victor.State) {
	warrantyCmd := exec.Command(w.Path, "--quit-on-error")

	// only accept 1 serial number in chat
	if len(s.Fields()) > 1 {
		msg := "Please only input one serial number at a time."
		s.Chat().Send(s.Message().Channel().ID(), msg)
		return
	}

	for _, arg := range s.Fields() {
		warrantyCmd.Args = append(warrantyCmd.Args, arg)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	warrantyCmd.Stdout = &stdout
	warrantyCmd.Stderr = &stderr
	err := warrantyCmd.Run()
	if err != nil {
		log.Println(err)
	}

	// Combine output and wrap stdout in code tags.
	output := fmt.Sprintf("```%s```\n%s", stdout.String(), stderr.String())

	// Send output and stderr to chat.
	s.Chat().Send(s.Message().Channel().ID(), output)
}

func init() {
	plugins.Add("warranty", func() plugins.Registrator {
		return &Warranty{}
	})
}
