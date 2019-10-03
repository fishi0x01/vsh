package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"io"
)

// RemoveCommand container for all 'rm' parameters
type RemoveCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Path   string
}

// NewRemoveCommand creates a new RemoveCommand parameter container
func NewRemoveCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *RemoveCommand {
	return &RemoveCommand{
		name:   "rm",
		client: c,
		stdout: stdout,
		stderr: stderr,
	}
}

// GetName returns the RemoveCommand's name identifier
func (cmd *RemoveCommand) GetName() string {
	return cmd.name
}

func (cmd *RemoveCommand) validate() error {
	log.Warn("Missing implementation of 'rm' validation")
	return nil
}

// Run executes 'rm' with given RemoveCommand's parameters
func (cmd *RemoveCommand) Run() error {
	err := cmd.validate()
	if err != nil {
		return err
	}

	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)

	for _, path := range cmd.client.Traverse(newPwd) {
		err := cmd.removeSecret(path)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.stdout, "Removed "+path)
	}

	return nil
}

func (cmd *RemoveCommand) removeSecret(path string) error {
	// delete
	err := cmd.client.Delete(path)
	if err != nil {
		return err
	}

	return nil
}
