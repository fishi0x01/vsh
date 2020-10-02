package cli

import (
	"fmt"
	"io"
	"strings"

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
func (cmd *CdCommand) Parse(args []string) error {
	if len(args) != 2 {
		fmt.Println("Usage:\ncd <path>")
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Path = args[1]
	return nil
}

// Run executes 'cd' with given CdCommand's parameters
func (cmd *CdCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)

	t := cmd.client.GetType(newPwd)

	if t == client.NONE {
		log.NotAValidPath(newPwd)
		return 1
	}

	if t == client.LEAF {
		log.NotAValidPath(newPwd)
		return 1
	}

	if !strings.HasSuffix(newPwd, "/") {
		newPwd = newPwd + "/"
	}
	cmd.client.Pwd = newPwd
	return 0
}
