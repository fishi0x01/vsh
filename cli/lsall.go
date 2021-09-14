package cli

import (
	"fmt"
	"strings"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// ListCommand container for 'ls' parameters
type ListAllCommand struct {
	name string
	args *ListAllCommandArgs

	client *client.Client
}

// ListCommandArgs provides a struct for go-arg parsing
type ListAllCommandArgs struct {
	Path string `arg:"positional" help:"path to list all child paths of, defaults to current path"`
}

// Description provides detail on what the command does
func (ListAllCommandArgs) Description() string {
	return "lists all child paths from the a given parent path"
}

// NewListCommand creates a new ListCommand parameter container
func NewListAllCommand(c *client.Client) *ListAllCommand {
	return &ListAllCommand{
		name:   "lsa",
		client: c,
		args:   &ListAllCommandArgs{},
	}
}

// GetName returns the ListCommand's name identifier
func (cmd *ListAllCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *ListAllCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *ListAllCommand) IsSane() bool {
	return cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *ListAllCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *ListAllCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}
	if cmd.args.Path == "" {
		cmd.args.Path = cmd.client.Pwd
	}

	return nil
}

// Run executes 'ls' with given ListAllCommand's parameters
func (cmd *ListAllCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.args.Path)
	result, err := cmd.client.ListAll(newPwd)

	if err != nil {
		log.UserError("Not a valid path for operation: %s", newPwd)
		return 1
	}
	log.UserInfo(strings.Join(result, "\n"))
	return 0
}
