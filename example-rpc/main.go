package main

import (
	"log"
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

func main() {
	log.SetPrefix("[plugin log] ")

	p := pie.NewProvider()
	if err := p.RegisterName("HelloRPC", api{}); err != nil {
		log.Fatalf("failed to register Plugin: %s", err)
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
}

type Args struct {
	User   string
	Fields []string
}

type api struct{}

func (api) Name(_ string, response *string) error {
	*response = "hello"
	return nil
}

func (api) Description(_ string, response *string) error {
	*response = "hello"
	return nil
}

func (api) Usage(_ string, response *[]string) error {
	*response = []string{"foo", "bar", "baz"}
	return nil
}

func (api) Handle(args *Args, response *string) error {
	*response = "Hi " + args.User
	return nil
}
