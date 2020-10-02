package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// RemoveCommand container for all 'rm' parameters
type RemoveCommand struct {
	name string

	client *client.Client
	Path   string
}

// NewRemoveCommand creates a new RemoveCommand parameter container
func NewRemoveCommand(c *client.Client) *RemoveCommand {
	return &RemoveCommand{
		name:   "rm",
		client: c,
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

// PrintUsage print command usage
func (cmd *RemoveCommand) PrintUsage() {
	log.UserInfo("Usage:\nrm <path>")
}

// Parse given arguments and return status
func (cmd *RemoveCommand) Parse(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Path = args[1]
	return nil
}

// Run executes 'rm' with given RemoveCommand's parameters
func (cmd *RemoveCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)

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
