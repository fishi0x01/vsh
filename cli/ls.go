package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
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

// IsSane returns true if command is sane
func (cmd *ListCommand) IsSane() bool {
	return cmd.Path != ""
}

// Parse given arguments and return status
func (cmd *ListCommand) Parse(args []string) (success bool) {
	if len(args) == 2 {
		cmd.Path = args[1]
		success = true
	} else if len(args) == 1 {
		cmd.Path = cmd.client.Pwd
		success = true
	} else {
		fmt.Println("Usage:\nls <path // optional>")
	}
	return success
}

// Run executes 'ls' with given ListCommand's parameters
func (cmd *ListCommand) Run() {
	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)
	result, err := cmd.client.List(newPwd)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintln(cmd.stdout, strings.Join(result, "\n"))
	}

	return
}
