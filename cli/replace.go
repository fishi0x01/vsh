package cli

import (
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// ReplaceMode defines replacement behaviour
type ReplaceMode string

const (
	// ModeValue only replace values
	ModeValue ReplaceMode = "value"
	// ModeKey only replace keys
	ModeKey ReplaceMode = "key"
	// ModeAll replace keys and values
	ModeAll AppendMode = "all"
)

// ReplaceCommand container for all 'replace' parameters
type ReplaceCommand struct {
	name string

	client   *client.Client
	Original string
	Replace  string
	Target   string
	Mode     ReplaceMode
}

// NewReplaceCommand creates a new ReplaceCommand parameter container
func NewReplaceCommand(c *client.Client) *ReplaceCommand {
	return &ReplaceCommand{
		name:   "replace",
		client: c,
		Mode:   ModeValue,
	}
}

// GetName returns the ReplaceCommand's name identifier
func (cmd *ReplaceCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *ReplaceCommand) IsSane() bool {
	return cmd.Original != "" && cmd.Target != ""
}

// PrintUsage print command usage
func (cmd *ReplaceCommand) PrintUsage() {
	log.UserInfo("Usage:\nreplace [--value|--key] <original> <replace> <file>")
}

// Parse parses the arguments and returns true on success; otherwise it prints usage and returns false
func (cmd *ReplaceCommand) Parse(args []string) error {
	// TODO
	return nil
}

// Run executes 'replace' with given ReplaceCommand's parameters
func (cmd *ReplaceCommand) Run() int {
	// TODO
	return 0
}
