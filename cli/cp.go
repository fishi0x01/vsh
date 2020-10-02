package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CopyCommand container for all 'cp' parameters
type CopyCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Source string
	Target string
}

// NewCopyCommand creates a new CopyCommand parameter container
func NewCopyCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *CopyCommand {
	return &CopyCommand{
		name:   "cp",
		client: c,
		stdout: stdout,
		stderr: stderr,
	}
}

// GetName returns the CopyCommand's name identifier
func (cmd *CopyCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *CopyCommand) IsSane() bool {
	return cmd.Source != "" && cmd.Target != ""
}

// Parse given arguments and return status
func (cmd *CopyCommand) Parse(args []string) error {
	if len(args) != 3 {
		fmt.Println("Usage:\ncp <from> <to>")
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Source = args[1]
	cmd.Target = args[2]
	return nil
}

// Run executes 'cp' with given CopyCommand's parameters
func (cmd *CopyCommand) Run() int {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	switch t := cmd.client.GetType(newSrcPwd); t {
	case client.LEAF:
		cmd.copySecret(filepath.Clean(newSrcPwd), newTargetPwd)
	case client.NODE:
		runCommandWithTraverseTwoPaths(cmd.client, newSrcPwd, newTargetPwd, cmd.copySecret)
	default:
		log.NotAValidPath(newSrcPwd)
		return 1
	}

	return 0
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

	log.Info("Copied %s to %s", source, target)

	return nil
}
