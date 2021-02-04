package cli

import (
	"fmt"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// MoveCommand container for all 'mv' parameters
type MoveCommand struct {
	name string
	args *MoveCommandArgs

	client *client.Client
}

// MoveCommandArgs provides a struct for go-arg parsing
type MoveCommandArgs struct {
	Source string `arg:"positional,required" help:"path to move"`
	Target string `arg:"positional,required" help:"path to move source to"`
}

// Description provides detail on what the command does
func (MoveCommandArgs) Description() string {
	return "moves a secret from one path to another"
}

// NewMoveCommand creates a new MoveCommand parameter container
func NewMoveCommand(c *client.Client) *MoveCommand {
	return &MoveCommand{
		name:   "mv",
		client: c,
		args:   &MoveCommandArgs{},
	}
}

// GetName returns the MoveCommand's name identifier
func (cmd *MoveCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *MoveCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *MoveCommand) IsSane() bool {
	return cmd.args.Source != "" && cmd.args.Target != ""
}

// PrintUsage print command usage
func (cmd *MoveCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *MoveCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Run executes 'mv' with given MoveCommand's parameters
func (cmd *MoveCommand) Run() int {
	return transportSecrets(cmd.client, cmd.args.Source, cmd.args.Target, cmd.moveSecret)
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
