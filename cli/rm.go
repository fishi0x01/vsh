package cli

import (
	"fmt"
	"io"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
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

// IsSane returns true if command is sane
func (cmd *RemoveCommand) IsSane() bool {
	return cmd.Path != ""
}

// Parse given arguments and return status
func (cmd *RemoveCommand) Parse(args []string) (success bool) {
	if len(args) == 2 {
		cmd.Path = args[1]
		success = true
	} else {
		fmt.Println("Usage:\nrm <path>")
	}
	return success
}

// Run executes 'rm' with given RemoveCommand's parameters
func (cmd *RemoveCommand) Run() {
	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)

	switch t := cmd.client.GetType(newPwd); t {
	case client.LEAF:
		cmd.removeSecret(newPwd)
	case client.NODE:
		for _, path := range cmd.client.Traverse(newPwd) {
			err := cmd.removeSecret(path)
			if err != nil {
				return
			}
		}
	default:
		log.Error("Invalid path: %s", newPwd)
	}
}

func (cmd *RemoveCommand) removeSecret(path string) error {
	// delete
	err := cmd.client.Delete(path)
	if err != nil {
		return err
	}

	log.Info("Removed %s", path)

	return nil
}
