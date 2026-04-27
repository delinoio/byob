package cli

import (
	"fmt"
	"io"

	bridge "github.com/microsoft/typescript-go/byobbridge"
)

const byobVersion = "0.0.0-dev"

type commandIdentifier string

const (
	commandHelp    commandIdentifier = "help"
	commandVersion commandIdentifier = "version"
	commandLint    commandIdentifier = "lint"
)

func Execute(args []string) int {
	return executeWithContext(args, defaultCommandContext())
}

func execute(args []string, stdout io.Writer, stderr io.Writer) int {
	ctx := defaultCommandContext()
	ctx.stdout = stdout
	ctx.stderr = stderr
	return executeWithContext(args, ctx)
}

func executeWithContext(args []string, ctx commandContext) int {
	ctx.withDefaults()

	if len(args) == 0 {
		printUsage(ctx.stderr)
		return 2
	}

	switch commandIdentifier(args[0]) {
	case "-h", "--help", commandHelp:
		printUsage(ctx.stdout)
		return 0
	case commandVersion:
		return executeVersion(args[1:], ctx.stdout, ctx.stderr)
	case commandLint:
		return executeLint(args[1:], ctx)
	default:
		_, _ = fmt.Fprintf(ctx.stderr, "unknown command %q\n", args[0])
		printUsage(ctx.stderr)
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
	_, _ = fmt.Fprintln(w, "  lint     Build and run BYOB lint tools")
}
