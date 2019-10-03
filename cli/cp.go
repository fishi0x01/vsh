package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"io"
)

// CopyCommand container for all 'cp' parameters
type CopyCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Source string
	Target string
}

// NewCopyCommand creates a new CopyCommand parameter container
func NewCopyCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *CopyCommand {
	return &CopyCommand{
		name:   "cp",
		client: c,
		stdout: stdout,
		stderr: stderr,
	}
}

// GetName returns the CopyCommand's name identifier
func (cmd *CopyCommand) GetName() string {
	return cmd.name
}

// Run executes 'cp' with given CopyCommand's parameters
func (cmd *CopyCommand) Run() error {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	return runCommandWithTraverseTwoPaths(cmd.client, newSrcPwd, newTargetPwd, cmd.copySecret)
}

func (cmd *CopyCommand) copySecret(source string, target string) error {
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

	fmt.Fprintln(cmd.stdout, "Moved "+source+" to "+target)

	return nil
}
