package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// MoveCommand container for all 'mv' parameters
type MoveCommand struct {
	name string

	client *client.Client
	Source string
	Target string
}

// NewMoveCommand creates a new MoveCommand parameter container
func NewMoveCommand(c *client.Client) *MoveCommand {
	return &MoveCommand{
		name:   "mv",
		client: c,
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

// PrintUsage print command usage
func (cmd *MoveCommand) PrintUsage() {
	log.UserInfo("Usage:\nmv <from> <to>")
}

// Parse given arguments and return status
func (cmd *MoveCommand) Parse(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Source = args[1]
	cmd.Target = args[2]
	return nil
}

// Run executes 'mv' with given MoveCommand's parameters
func (cmd *MoveCommand) Run() int {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	switch t := cmd.client.GetType(newSrcPwd); t {
	case client.LEAF:
		cmd.moveSecret(filepath.Clean(newSrcPwd), newTargetPwd)
	case client.NODE:
		runCommandWithTraverseTwoPaths(cmd.client, newSrcPwd, newTargetPwd, cmd.moveSecret)
	default:
		log.UserError("Not a valid path for operation: %s", newSrcPwd)
		return 1
	}

	return 0
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

	log.UserDebug("Moved %s to %s", source, target)

	return nil
}
