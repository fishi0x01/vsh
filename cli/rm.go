package cli

import (
	"fmt"
	"sync"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// RemoveCommand container for all 'rm' parameters
type RemoveCommand struct {
	name        string
	args        *RemoveCommandArgs
	workerCount int

	client *client.Client
}

// RemoveCommandArgs provides a struct for go-arg parsing
type RemoveCommandArgs struct {
	Recursive bool   `arg:"-r" help:"recursively remove a directory"`
	Path      string `arg:"positional,required" help:"path to remove"`
}

// Description provides detail on what the command does
func (RemoveCommandArgs) Description() string {
	return "removes a secret at a path"
}

// NewRemoveCommand creates a new RemoveCommand parameter container
func NewRemoveCommand(c *client.Client, workerCount int) *RemoveCommand {
	return &RemoveCommand{
		name:        "rm",
		client:      c,
		args:        &RemoveCommandArgs{},
		workerCount: workerCount,
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
			fmt.Printf("Error removing secret: %v", err)
			return 1
		}
	case client.NODE:
		if !cmd.args.Recursive {
			log.UserError("use -r to remove directories")
			return 1
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
					fmt.Printf("Error removing dir: %v", err)
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
		log.UserError("not a valid path for operation: %s", newPwd)
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
