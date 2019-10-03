package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"io"
	"strings"
)

// ListCommand container for 'ls' parameters
type ListCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Path   string
}

// NewListCommand creates a new ListCommand parameter container
func NewListCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *ListCommand {
	return &ListCommand{
		name:   "ls",
		client: c,
		stdout: stdout,
		stderr: stderr,
	}
}

// GetName returns the ListCommand's name identifier
func (cmd *ListCommand) GetName() string {
	return cmd.name
}

func (cmd *ListCommand) validate() error {
	log.Warn("Missing implementation of 'ls' validation")
	return nil
}

// Run executes 'ls' with given ListCommand's parameters
func (cmd *ListCommand) Run() error {
	err := cmd.validate()
	if err != nil {
		return err
	}

	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)
	result, err := cmd.client.List(newPwd)

	fmt.Fprintln(cmd.stdout, strings.Join(result, "  "))

	return err
}
