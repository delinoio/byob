package cli

import (
	"io"
	"os"
)

type commandContext struct {
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer
	cacheRoot string
	env       []string
}

func defaultCommandContext() commandContext {
	return commandContext{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		env:    os.Environ(),
	}
}

func (ctx *commandContext) withDefaults() {
	if ctx.stdin == nil {
		ctx.stdin = os.Stdin
	}
	if ctx.stdout == nil {
		ctx.stdout = io.Discard
	}
	if ctx.stderr == nil {
		ctx.stderr = io.Discard
	}
	if ctx.env == nil {
		ctx.env = os.Environ()
	}
}
