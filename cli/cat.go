package cli

import (
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// CatCommand container for all 'cat' parameters
type CatCommand struct {
	name string

	client *client.Client
	Path   string
}

// NewCatCommand creates a new CatCommand parameter container
func NewCatCommand(c *client.Client) *CatCommand {
	return &CatCommand{
		name:   "cat",
		client: c,
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

// PrintUsage print command usage
func (cmd *CatCommand) PrintUsage() {
	log.UserInfo("Usage:\ncat <secret>")
}

// Parse given arguments and return status
func (cmd *CatCommand) Parse(args []string) error {
	if len(args) != 2 {
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
					log.UserInfo("%s = %s", kk, vv)
				}
			} else {
				// KV 1
				log.UserInfo("%s = %s", k, v)
			}
		}
	} else {
		log.UserError("Not a valid path for operation: %s", absPath)
		return 1
	}
	return 0
}
