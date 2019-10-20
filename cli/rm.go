package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
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

	t := cmd.client.GetType(newPwd)
	if t != client.NODE && t != client.LEAF {
		fmt.Fprintln(cmd.stderr, "Not a valid path: "+newPwd)
		return
	}

	for _, path := range cmd.client.Traverse(newPwd) {
		err := cmd.removeSecret(path)
		if err != nil {
			return
		}
	}

	return
}

func (cmd *RemoveCommand) removeSecret(path string) error {
	// delete
	err := cmd.client.Delete(path)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.stdout, "Removed "+path)

	return nil
}
