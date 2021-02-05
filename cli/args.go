package cli

import (
	"bytes"

	"github.com/alexflint/go-arg"
)

func parseCommandArgs(args []string, cmd Command) (*arg.Parser, error) {
	p, err := argParser(args, cmd)
	if err != nil {
		return nil, err
	}

	if len(args) > 1 {
		err = p.Parse(args[1:])
	} else {
		err = p.Parse([]string{})
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func argParser(args []string, cmd Command) (*arg.Parser, error) {
	return arg.NewParser(arg.Config{Program: args[0]}, cmd.GetArgs())
}

// Usage returns usage information
func Usage(cmd Command) string {
	var b bytes.Buffer
	p, _ := argParser([]string{cmd.GetName()}, cmd)
	p.WriteUsage(&b)
	return b.String()
}

// Help returns extended usage information
func Help(cmd Command) string {
	var b bytes.Buffer
	p, _ := argParser([]string{cmd.GetName()}, cmd)
	p.WriteHelp(&b)
	return b.String()
}
