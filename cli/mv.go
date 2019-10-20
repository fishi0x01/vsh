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

// IsSane returns true if command is sane
func (cmd *MoveCommand) IsSane() bool {
	return cmd.Source != "" && cmd.Target != ""
}

// Parse given arguments and return status
func (cmd *MoveCommand) Parse(args []string) (success bool) {
	if len(args) == 3 {
		cmd.Source = args[1]
		cmd.Target = args[2]
		success = true
	} else {
		fmt.Println("Usage:\nmv <from> <to>")
	}
	return success
}

// Run executes 'mv' with given MoveCommand's parameters
func (cmd *MoveCommand) Run() {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	t := cmd.client.GetType(newSrcPwd)
	if t != client.NODE && t != client.LEAF {
		fmt.Fprintln(cmd.stderr, "Not a valid source path: "+newSrcPwd)
		return
	}

	runCommandWithTraverseTwoPaths(cmd.client, newSrcPwd, newTargetPwd, cmd.moveSecret)
	return
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
