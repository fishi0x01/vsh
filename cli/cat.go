package cli

import (
	"fmt"
	"io"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CatCommand container for all 'cat' parameters
type CatCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Path   string
}

// NewCatCommand creates a new CatCommand parameter container
func NewCatCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *CatCommand {
	return &CatCommand{
		name:   "cat",
		client: c,
		stderr: stderr,
		stdout: stdout,
	}
}

// GetName returns the CatCommand's name identifier
func (cmd *CatCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *CatCommand) IsSane() bool {
	return cmd.Path != ""
}

// Parse given arguments and return status
func (cmd *CatCommand) Parse(args []string) error {
	if len(args) != 2 {
		fmt.Println("Usage:\ncat <secret>")
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Path = args[1]
	return nil
}

// Run executes 'cat' with given CatCommand's parameters
func (cmd *CatCommand) Run() int {
	absPath := cmdPath(cmd.client.Pwd, cmd.Path)
	t := cmd.client.GetType(absPath)

	if t == client.LEAF {
		secret, err := cmd.client.Read(absPath)
		if err != nil {
			return 1
		}

		for k, v := range secret.Data {
			if rec, ok := v.(map[string]interface{}); ok {
				// KV 2
				for kk, vv := range rec {
					fmt.Fprintln(cmd.stdout, kk, "=", vv)
				}
			} else {
				// KV 1
				fmt.Fprintln(cmd.stdout, k, "=", v)
			}
		}
	} else {
		log.NotAValidPath(absPath)
		return 1
	}
	return 0
}
