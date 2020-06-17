package cli

import (
	"fmt"
	"io"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CdCommand container for all 'cd' parameters
type CdCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Path   string
}

// NewCdCommand creates a new CdCommand parameter container
func NewCdCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *CdCommand {
	return &CdCommand{
		name:   "cd",
		client: c,
		stdout: stdout,
		stderr: stderr,
	}
}

// GetName returns the CdCommand's name identifier
func (cmd *CdCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *CdCommand) IsSane() bool {
	return cmd.Path != ""
}

// Parse given arguments and return status
func (cmd *CdCommand) Parse(args []string) (success bool) {
	if len(args) == 2 {
		cmd.Path = args[1]
		success = true
	} else {
		fmt.Println("Usage:\ncd <path>")
	}
	return success
}

// Run executes 'cd' with given CdCommand's parameters
func (cmd *CdCommand) Run() {
	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)

	t := cmd.client.GetType(newPwd)

	if t == client.NONE {
		log.Error("Invalid directory: %s", newPwd)
		return
	}

	if t == client.LEAF {
		log.Error("Invalid directory: %s is a file", newPwd)
		return
	}

	newPwd = newPwd + "/"
	if newPwd == "//" {
		newPwd = "/"
	}
	cmd.client.Pwd = newPwd
	return
}
