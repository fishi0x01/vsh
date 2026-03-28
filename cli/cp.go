package cli

import (
	"fmt"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CopyCommand container for all 'cp' parameters
type CopyCommand struct {
	name        string
	args        *CopyCommandArgs
	workerCount int

	client *client.Client
}

// CopyCommandArgs provides a struct for go-arg parsing
type CopyCommandArgs struct {
	Recursive bool   `arg:"-r"                  help:"recursively copy a directory"`
	Source    string `arg:"positional,required" help:"path to copy from"`
	Target    string `arg:"positional,required" help:"path to copy to"`
}

// Description provides detail on what the command does
func (CopyCommandArgs) Description() string {
	return "recursively copies a path to another location"
}

// NewCopyCommand creates a new CopyCommand parameter container
func NewCopyCommand(c *client.Client, workerCount int) *CopyCommand {
	return &CopyCommand{
		name:        "cp",
		client:      c,
		args:        &CopyCommandArgs{},
		workerCount: workerCount,
	}
}

// GetName returns the CopyCommand's name identifier
func (cmd *CopyCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *CopyCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *CopyCommand) IsSane() bool {
	return cmd.args.Source != "" && cmd.args.Target != ""
}

// PrintUsage print command usage
func (cmd *CopyCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *CopyCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

// Run executes 'cp' with given CopyCommand's parameters
func (cmd *CopyCommand) Run() int {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.args.Source)
	if cmd.client.GetType(newSrcPwd) == client.NODE && !cmd.args.Recursive {
		log.UserError("use -r to copy directories")
		return 1
	}
	return transportSecrets(
		cmd.client, cmd.args.Source, cmd.args.Target, cmd.copySecret, cmd.workerCount,
	)
}

func (cmd *CopyCommand) copySecret(source string, target string) error {
	// read
	secret, err := cmd.client.Read(source)
	if err != nil {
		return err
	}

	// write
	err = cmd.client.Write(target, secret)
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.UserDebug("Copied %s to %s", source, target)

	return nil
}
