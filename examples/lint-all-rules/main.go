package main

import (
	"os"

	"github.com/delinoio/byob/lint/upstream"
)

func main() {
	os.Exit(upstream.RunAllRulesCLI(os.Args[1:]))
}
