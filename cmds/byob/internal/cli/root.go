package cli

import (
	"fmt"
	"io"
	"os"

	bridge "github.com/microsoft/typescript-go/byobbridge"
)

const byobVersion = "0.0.0-dev"

func Execute(args []string) int {
	return execute(args, os.Stdout, os.Stderr)
}

func execute(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage(stdout)
		return 0
	case "version":
		return executeVersion(args[1:], stdout, stderr)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func executeVersion(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(stderr, "version command accepts no arguments")
		printUsage(stderr)
		return 2
	}

	info := bridge.Info()
	_, _ = fmt.Fprintf(stdout, "byob version: %s\n", byobVersion)
	_, _ = fmt.Fprintf(stdout, "typescript-go module: %s\n", info.Module)
	_, _ = fmt.Fprintf(stdout, "typescript-go version: %s\n", info.Version)
	return 0
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "usage: byob <command>")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "commands:")
	_, _ = fmt.Fprintln(w, "  version  Print BYOB and TypeScript-Go versions")
}
