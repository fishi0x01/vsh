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

	absPath := cmd.client.Pwd + cmd.Path
	isFile, err := cmd.client.IsFile(absPath)
	if err != nil {
		return err
	}

	if isFile {
		secret, err := cmd.client.Read(absPath)
		if err != nil {
			return err
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
		fmt.Fprintln(cmd.stderr, cmd.client.Pwd+cmd.Path, "is not a file")
	}

	return err
}
