package main

import (
	"github.com/spf13/cobra"

	apiserver "ccrayz/event-scanner/cmd/api-server"
)

var command *cobra.Command

func main() {
	command = apiserver.NewCommand()

	if err := command.Execute(); err != nil {
		panic(err)
	}
}
