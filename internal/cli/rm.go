package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fishi0x01/vsh/internal/client"
	"github.com/fishi0x01/vsh/internal/logger"
)

// RemoveCommand container for all 'rm' parameters
type RemoveCommand struct {
	name        string
	args        *RemoveCommandArgs
	workerCount int
	interactive bool

	client *client.Client
}

// RemoveCommandArgs provides a struct for go-arg parsing
type RemoveCommandArgs struct {
	Recursive bool   `arg:"-r"                  help:"recursively remove a directory"`
	Force     bool   `arg:"-f"                  help:"skip confirmation prompt for recursive removal"`
	Path      string `arg:"positional,required" help:"path to remove"`
}

// Description provides detail on what the command does
func (RemoveCommandArgs) Description() string {
	return "removes a secret at a path"
}

// NewRemoveCommand creates a new RemoveCommand parameter container
func NewRemoveCommand(c *client.Client, workerCount int, interactive bool) *RemoveCommand {
	return &RemoveCommand{
		name:        "rm",
		client:      c,
		args:        &RemoveCommandArgs{},
		workerCount: workerCount,
		interactive: interactive,
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
		err := cmd.removeSecret(newPwd)
		if err != nil {
			logger.UserError("Error removing secret: %v", err)
			return 1
		}
	case client.NODE:
		if !cmd.args.Recursive {
			logger.UserError("use -r to remove directories")
			return 1
		}
		if cmd.interactive && !cmd.args.Force {
			fmt.Printf("Remove entire subtree at '%s'? [y/N] ", newPwd)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(answer)) != "y" {
				fmt.Println("Aborted.")
				return 0
			}
		}
		var wg sync.WaitGroup
		sem := make(chan struct{}, cmd.workerCount)
		failed := make(chan struct{}, 1)
		for _, path := range cmd.client.Traverse(newPwd, false) {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				if err := cmd.removeSecret(p); err != nil {
					logger.UserError("Error removing dir: %v", err)
					select {
					case failed <- struct{}{}:
					default:
					}
				}
			}(path)
		}
		wg.Wait()
		if len(failed) > 0 {
			return 1
		}
	default:
		logger.UserError("not a valid path for operation: %s", newPwd)
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

	logger.UserDebug("Removed %s", path)

	return nil
}
