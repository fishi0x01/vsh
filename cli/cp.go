package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"io"
	"strings"
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

func (cmd *CopyCommand) validate() error {
	log.Warn("Missing implementation of 'cp' validation")
	return nil
}

// Run executes 'cp' with given MoveCommand's parameters
func (cmd *CopyCommand) Run() error {
	err := cmd.validate()

	if err != nil {
		return err
	}

	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	for _, path := range cmd.client.Traverse(newSrcPwd) {
		target := strings.Replace(path, newSrcPwd, newTargetPwd, 1)
		err := cmd.copySecret(path, target)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.stdout, "Moved "+path+" to "+target)
	}

	return nil
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

	return nil
}
