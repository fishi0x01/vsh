package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"io"
	"strings"
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

func (cmd *MoveCommand) validate() error {
	log.Warn("Missing implementation of 'mv' validation")
	return nil
}

// Run executes 'mv' with given MoveCommand's parameters
func (cmd *MoveCommand) Run() error {
	err := cmd.validate()

	if err != nil {
		return err
	}

	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	for _, path := range cmd.client.Traverse(newSrcPwd) {
		target := strings.Replace(path, newSrcPwd, newTargetPwd, 1)
		err := cmd.moveSecret(path, target)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.stdout, "Moved "+path+" to "+target)
	}

	return nil
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

	return nil
}
