package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"io"
)

// MoveCommand container for all 'mv' parameters
type MoveCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Source string
	Target string
}

// NewMoveCommand creates a new MoveCommand parameter container
func NewMoveCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *MoveCommand {
	return &MoveCommand{
		name:   "mv",
		client: c,
		stdout: stdout,
		stderr: stderr,
	}
}

// GetName returns the MoveCommand's name identifier
func (cmd *MoveCommand) GetName() string {
	return cmd.name
}

// Run executes 'mv' with given MoveCommand's parameters
func (cmd *MoveCommand) Run() error {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	return runCommandWithTraverseTwoPaths(cmd.client, newSrcPwd, newTargetPwd, cmd.moveSecret)
}

func (cmd *MoveCommand) moveSecret(source string, target string) error {
	// read
	secret, err := cmd.client.Read(source)
	if err != nil {
		return err
	}

	// write
	err = cmd.client.Write(target, secret)
	if err != nil {
		return err
	}

	// delete
	err = cmd.client.Delete(source)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.stdout, "Moved "+source+" to "+target)

	return nil
}
