package main

import (
	logger "github.com/hashicorp/go-hclog"

	"github.com/umbracle/minimal/command"
	_ "github.com/umbracle/minimal/command/version"
	_ "github.com/umbracle/minimal/command/peers"
	_ "github.com/umbracle/minimal/command/genesis"
	_ "github.com/umbracle/minimal/command/agent"
	_ "github.com/umbracle/minimal/command/dev"
	_ "github.com/umbracle/minimal/command/debug"
)

func main() {
	// TODO: Change time format for the logger?
	if err := command.Run(); err != nil {
		logger.Default().Error(err.Error())
	}
}
