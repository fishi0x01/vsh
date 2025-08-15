package cli

import (
	"fmt"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CatCommand container for all 'cat' parameters
type CatCommand struct {
	name string
	args *CatCommandArgs

	client *client.Client
}

// CatCommandArgs provides a struct for go-arg parsing
type CatCommandArgs struct {
	Path string `arg:"positional,required" help:"path to display contents"`
}

// Description provides detail on what the command does
func (CatCommandArgs) Description() string {
	return "displays the content of a secret"
}

// NewCatCommand creates a new CatCommand parameter container
func NewCatCommand(c *client.Client) *CatCommand {
	return &CatCommand{
		name:   "cat",
		client: c,
		args:   &CatCommandArgs{},
	}
}

// GetName returns the CatCommand's name identifier
func (cmd *CatCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *CatCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *CatCommand) IsSane() bool {
	return cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *CatCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *CatCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Run executes 'cat' with given CatCommand's parameters
func (cmd *CatCommand) Run() int {
	absPath := cmdPath(cmd.client.Pwd, cmd.args.Path)
	t := cmd.client.GetType(absPath)

	if t == client.LEAF {
		secret, err := cmd.client.Read(absPath)
		if err != nil {
			return 1
		}

		for k, v := range secret.GetData() {
			log.UserInfo("%s = %s", k, v)
		}
	} else {
		log.UserError("not a valid path for operation: %s", absPath)
		return 1
	}
	return 0
}
