package cli

import (
	"fmt"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// RemoveCommand container for all 'rm' parameters
type RemoveCommand struct {
	name string
	args *RemoveCommandArgs

	client *client.Client
}

// RemoveCommandArgs provides a struct for go-arg parsing
type RemoveCommandArgs struct {
	Path string `arg:"positional,required" help:"path to remove"`
}

// Description provides detail on what the command does
func (RemoveCommandArgs) Description() string {
	return "removes a secret at a path"
}

// NewRemoveCommand creates a new RemoveCommand parameter container
func NewRemoveCommand(c *client.Client) *RemoveCommand {
	return &RemoveCommand{
		name:   "rm",
		client: c,
		args:   &RemoveCommandArgs{},
	}
}

// GetName returns the RemoveCommand's name identifier
func (cmd *RemoveCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *RemoveCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *RemoveCommand) IsSane() bool {
	return cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *RemoveCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *RemoveCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Run executes 'rm' with given RemoveCommand's parameters
func (cmd *RemoveCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.args.Path)

	switch t := cmd.client.GetType(newPwd); t {
	case client.LEAF:
		cmd.removeSecret(newPwd)
	case client.NODE:
		for _, path := range cmd.client.Traverse(newPwd) {
			err := cmd.removeSecret(path)
			if err != nil {
				return 1
			}
		}
	default:
		log.UserError("Not a valid path for operation: %s", newPwd)
		return 1
	}

	return 0
}

func (cmd *RemoveCommand) removeSecret(path string) error {
	// delete
	err := cmd.client.Delete(path)
	if err != nil {
		return err
	}

	log.UserDebug("Removed %s", path)

	return nil
}
