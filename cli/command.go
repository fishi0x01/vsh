package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fishi0x01/vsh/client"
)

// Command interface to describe a command structure
type Command interface {
	Run() int
	GetName() string
	IsSane() bool
	PrintUsage()
	Parse(args []string) error
}

// Commands contains all available commands
type Commands struct {
	Mv      *MoveCommand
	Cp      *CopyCommand
	Append  *AppendCommand
	Rm      *RemoveCommand
	Ls      *ListCommand
	Cd      *CdCommand
	Cat     *CatCommand
	Grep    *GrepCommand
	Replace *ReplaceCommand
}

// NewCommands returns a Commands struct with all available commands
func NewCommands(client *client.Client) *Commands {
	return &Commands{
		Mv:      NewMoveCommand(client),
		Cp:      NewCopyCommand(client),
		Append:  NewAppendCommand(client),
		Rm:      NewRemoveCommand(client),
		Ls:      NewListCommand(client),
		Cd:      NewCdCommand(client),
		Cat:     NewCatCommand(client),
		Grep:    NewGrepCommand(client, os.Stdout, os.Stderr),
		Replace: NewReplaceCommand(client),
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
	for _, path := range client.Traverse(source) {
		target := strings.Replace(path, source, target, 1)
		err := f(path, target)
		if err != nil {
			return
		}
	}

	return
}
