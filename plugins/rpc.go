package plugins

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/FogCreek/victor"
	"github.com/natefinch/pie"
)

type RPCPlugin struct {
	name string
	path string

	cmdName        string
	cmdDescription string
	cmdUsage       []string

	handlers *[]victor.HandlerDocPair

	client *rpc.Client
}

func NewRPCPlugin(name, path string) *RPCPlugin {
	plugin := &RPCPlugin{name: name, path: path + "/" + name}
	err := plugin.newClient()
	if err != nil {
		log.Fatal(err)
	}
	defer plugin.client.Close()

	err = plugin.setCmdName()
	if err != nil {
		log.Fatal(err)
	}
	err = plugin.setCmdDescription()
	if err != nil {
		log.Fatal(err)
	}
	err = plugin.setCmdUsage()
	if err != nil {
		log.Fatal(err)
	}
	plugin.createHandlers()

	return plugin
}

func (p *RPCPlugin) newClient() error {
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, p.path)
	if err != nil {
		log.Fatalf("Error running plugin: %s", err)
	}
	p.client = client
	return nil
}

func (p RPCPlugin) Register() []victor.HandlerDocPair {
	return *p.handlers
}

func (p *RPCPlugin) setCmdName() error {
	return p.client.Call(fmt.Sprintf("%v.Name", p.name), nil, &p.cmdName)
}

func (p *RPCPlugin) setCmdDescription() error {
	return p.client.Call(fmt.Sprintf("%v.Description", p.name), nil, &p.cmdDescription)
}

func (p *RPCPlugin) setCmdUsage() error {
	return p.client.Call(fmt.Sprintf("%v.Usage", p.name), nil, &p.cmdUsage)
}

func (p *RPCPlugin) createHandlers() {
	handlers := &[]victor.HandlerDocPair{
		&victor.HandlerDoc{
			CmdHandler:     p.handleFunc,
			CmdName:        p.cmdName,
			CmdDescription: p.cmdDescription,
			CmdUsage:       p.cmdUsage,
		},
	}
	p.handlers = handlers
}

func (p RPCPlugin) handleFunc(s victor.State) {
	type Args struct {
		User   string
		Fields []string
	}
	var msg string
	args := &Args{User: s.Message().User().Name(), Fields: s.Fields()}
	err := p.newClient()
	if err != nil {
		log.Fatal(err)
	}
	err = p.client.Call(fmt.Sprintf("%v.Handle", p.name),
		args, &msg)
	if err != nil {
		log.Println(err)
		msg = fmt.Sprintf("Plugin encountered an error, %v", err)
	}
	s.Reply(string(msg))
}
