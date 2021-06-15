package cli

import (
	"path/filepath"
	"strings"

	"github.com/fatih/structs"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// Command interface to describe a command structure
type Command interface {
	Run() int
	GetName() string
	GetArgs() interface{}
	IsSane() bool
	PrintUsage()
	Parse(args []string) error
}

// Commands contains all available commands
type Commands struct {
	Add     *AddCommand
	Append  *AppendCommand
	Cat     *CatCommand
	Cd      *CdCommand
	Cp      *CopyCommand
	Grep    *GrepCommand
	Ls      *ListCommand
	Mv      *MoveCommand
	Replace *ReplaceCommand
	Rm      *RemoveCommand
}

// Get returns the Command that matches the string
func (cmds *Commands) Get(cmd string) Command {
	for _, f := range structs.Fields(cmds) {
		if c := f.Value().(Command); cmd == c.GetName() {
			return c
		}
	}
	return nil
}

// NewCommands returns a Commands struct with all available commands
func NewCommands(client *client.Client) *Commands {
	return &Commands{
		Add:     NewAddCommand(client),
		Append:  NewAppendCommand(client),
		Cat:     NewCatCommand(client),
		Cd:      NewCdCommand(client),
		Cp:      NewCopyCommand(client),
		Grep:    NewGrepCommand(client),
		Ls:      NewListCommand(client),
		Mv:      NewMoveCommand(client),
		Replace: NewReplaceCommand(client),
		Rm:      NewRemoveCommand(client),
	}
}

func cmdPath(pwd string, arg string) (result string) {
	result = filepath.Clean(pwd + arg)

	if strings.HasSuffix(arg, "/") {
		// filepath.Clean removes "/" suffix, but we need it to distinguish path from file
		result = result + "/"
	}

	if strings.HasPrefix(arg, "/") {
		// absolute path is given
		result = arg
	}
	return result
}

func runCommandWithTraverseTwoPaths(client *client.Client, source string, target string, f func(string, string) error) {
	source = filepath.Clean(source) // remove potential trailing '/'
	for _, path := range client.Traverse(source, false) {
		target := strings.Replace(path, source, target, 1)
		err := f(path, target)
		if err != nil {
			return
		}
	}

	return
}

func transportSecrets(c *client.Client, source string, target string, transport func(string, string) error) int {
	newSrcPwd := cmdPath(c.Pwd, source)
	newTargetPwd := cmdPath(c.Pwd, target)

	switch t := c.GetType(newSrcPwd); t {
	case client.LEAF:
		transport(filepath.Clean(newSrcPwd), newTargetPwd)
	case client.NODE:
		runCommandWithTraverseTwoPaths(c, newSrcPwd, newTargetPwd, transport)
	default:
		log.UserError("Not a valid path for operation: %s", newSrcPwd)
		return 1
	}

	return 0
}
