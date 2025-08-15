package cli

import (
	"fmt"
	"strings"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CdCommand container for all 'cd' parameters
type CdCommand struct {
	name string
	args *CdCommandArgs

	client *client.Client
}

// CdCommandArgs provides a struct for go-arg parsing
type CdCommandArgs struct {
	Path string `arg:"positional,required" help:"change cwd to path"`
}

// Description provides detail on what the command does
func (CdCommandArgs) Description() string {
	return "changes the working path"
}

// NewCdCommand creates a new CdCommand parameter container
func NewCdCommand(c *client.Client) *CdCommand {
	return &CdCommand{
		name:   "cd",
		client: c,
		args:   &CdCommandArgs{},
	}
}

// GetName returns the CdCommand's name identifier
func (cmd *CdCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *CdCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *CdCommand) IsSane() bool {
	return cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *CdCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *CdCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Run executes 'cd' with given CdCommand's parameters
func (cmd *CdCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.args.Path)

	t := cmd.client.GetType(newPwd)

	if t == client.NONE {
		log.UserError("not a valid path for operation: %s", newPwd)
		return 1
	}

	if t == client.LEAF {
		log.UserError("not a valid path for operation: %s", newPwd)
		return 1
	}

	if !strings.HasSuffix(newPwd, "/") {
		newPwd = newPwd + "/"
	}
	cmd.client.Pwd = newPwd
	return 0
}
