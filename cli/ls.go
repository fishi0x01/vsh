package cli

import (
	"fmt"
	"strings"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// ListCommand container for 'ls' parameters
type ListCommand struct {
	name string
	args *ListCommandArgs

	client *client.Client
}

// ListCommandArgs provides a struct for go-arg parsing
type ListCommandArgs struct {
	Path string `arg:"positional" help:"path to list contents of, defaults to current path"`
}

// Description provides detail on what the command does
func (ListCommandArgs) Description() string {
	return "lists the secrets at a path"
}

// NewListCommand creates a new ListCommand parameter container
func NewListCommand(c *client.Client) *ListCommand {
	return &ListCommand{
		name:   "ls",
		client: c,
		args:   &ListCommandArgs{},
	}
}

// GetName returns the ListCommand's name identifier
func (cmd *ListCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *ListCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *ListCommand) IsSane() bool {
	return cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *ListCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *ListCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}
	if cmd.args.Path == "" {
		cmd.args.Path = cmd.client.Pwd
	}

	return nil
}

// Run executes 'ls' with given ListCommand's parameters
func (cmd *ListCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.args.Path)
	result, err := cmd.client.List(newPwd)

	if err != nil {
		log.UserError("not a valid path for operation: %s", newPwd)
		return 1
	}
	log.UserInfo("%s", strings.Join(result, "\n"))
	return 0
}
