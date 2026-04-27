package main

import (
	"os"

	"github.com/delinoio/byob/cmds/byob/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:]))
}
