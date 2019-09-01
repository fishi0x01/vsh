package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"io"
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

func (cmd *CatCommand) validate() error {
	log.Warn("Missing implementation of 'cat' validation")
	return nil
}

// Run executes 'cat' with given CatCommand's parameters
func (cmd *CatCommand) Run() error {
	err := cmd.validate()
	if err != nil {
		return err
	}

	isFile, err := cmd.client.IsFile(cmd.client.Pwd + cmd.Path)
	if err != nil {
		return err
	}

	if isFile {
		secret, err := cmd.client.Read(cmd.client.Pwd + cmd.Path)
		if err != nil {
			return err
		}

		for k, v := range secret.Data {
			fmt.Fprintln(cmd.stdout, k, "=", v)
		}
	} else {
		fmt.Fprintln(cmd.stderr, "'", cmd.client.Pwd+cmd.Path, "' is not a file")
	}

	return err
}
