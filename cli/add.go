package cli

import (
	"fmt"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// AddCommand container for all 'append' parameters
type AddCommand struct {
	name string
	args *AddCommandArgs

	client *client.Client
}

// AddCommandArgs provides a struct for go-arg parsing
type AddCommandArgs struct {
	Key   string `arg:"positional,required"`
	Value string `arg:"positional,required"`
	Path  string `arg:"positional,required"`
	Force bool   `arg:"-f,--force" help:"Overwrite key if exists"`
}

// Description provides detail on what the command does
func (AddCommandArgs) Description() string {
	return "adds a key with value to a path"
}

// NewAddCommand creates a new AddCommand parameter container
func NewAddCommand(c *client.Client) *AddCommand {
	return &AddCommand{
		name:   "add",
		client: c,
		args:   &AddCommandArgs{},
	}
}

// GetName returns the AddCommand's name identifier
func (cmd *AddCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *AddCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *AddCommand) IsSane() bool {
	return cmd.args.Key != "" && cmd.args.Value != "" && cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *AddCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse parses the arguments into the Command and Args structs
func (cmd *AddCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Run executes 'add' with given AddCommand's parameters
func (cmd *AddCommand) Run() int {
	path := cmdPath(cmd.client.Pwd, cmd.args.Path)

	pathType := cmd.client.GetType(path)
	if pathType != client.LEAF {
		log.UserError("Not a valid path for operation: %s", path)
		return 1
	}

	err := cmd.addKeyValue(cmd.args.Path, cmd.args.Key, cmd.args.Value)
	if err != nil {
		log.UserError("Add failed: " + err.Error())
		return 1
	}

	return 0
}

func (cmd *AddCommand) addKeyValue(path string, key string, value string) error {
	secret, err := cmd.client.Read(path)
	if err != nil {
		return fmt.Errorf("Read failed: %s", err)
	}
	data := secret.GetData()
	if _, ok := data[key]; ok && !cmd.args.Force {
		return fmt.Errorf("Key already exists at path: %s", path)
	}
	data[key] = value
	secret.SetData(data)
	return cmd.client.Write(path, secret)
}
