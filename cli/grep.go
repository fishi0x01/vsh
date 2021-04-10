package cli

import (
	"fmt"
	"os"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// GrepCommand container for all 'grep' parameters
type GrepCommand struct {
	name string
	args *GrepCommandArgs

	client   *client.Client
	searcher *Searcher
	Mode     KeyValueMode
}

// GrepCommandArgs provides a struct for go-arg parsing
type GrepCommandArgs struct {
	Search string `arg:"positional,required"`
	Path   string `arg:"positional"`
	Regexp bool   `arg:"-e,--regexp" help:"Treat search string as a regexp"`
	Keys   bool   `arg:"-k,--keys" help:"Match against keys (true if -v is not specified)"`
	Values bool   `arg:"-v,--values" help:"Match against values (true if -k is not specified)"`
}

// Description provides detail on what the command does
func (GrepCommandArgs) Description() string {
	return "recursive searches for a pattern starting at a path"
}

// NewGrepCommand creates a new GrepCommand parameter container
func NewGrepCommand(c *client.Client) *GrepCommand {
	return &GrepCommand{
		name:   "grep",
		client: c,
		args:   &GrepCommandArgs{},
	}
}

// GetName returns the GrepCommand's name identifier
func (cmd *GrepCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *GrepCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *GrepCommand) IsSane() bool {
	return cmd.args.Path != "" && cmd.args.Search != ""
}

// PrintUsage print command usage
func (cmd *GrepCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse given arguments and return status
func (cmd *GrepCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}
	if cmd.args.Path == "" {
		cmd.args.Path = cmd.client.Pwd
	}
	if cmd.args.Keys == true {
		cmd.Mode |= ModeKeys
	}
	if cmd.args.Values == true {
		cmd.Mode |= ModeValues
	}
	if cmd.Mode == 0 {
		cmd.Mode = ModeKeys + ModeValues
	}

	searcher, err := NewSearcher(cmd)
	if err != nil {
		return err
	}
	cmd.searcher = searcher

	return nil
}

// Run executes 'grep' with given GrepCommand's parameters
func (cmd *GrepCommand) Run() int {
	path := cmdPath(cmd.client.Pwd, cmd.args.Path)
	filePaths, err := cmd.client.SubpathsForPath(path)
	if err != nil {
		log.UserError(fmt.Sprintf("%s", err))
		return 1
	}

	for _, curPath := range filePaths {
		matches, err := cmd.grepFile(cmd.args.Search, curPath)
		if err != nil {
			return 1
		}
		for _, match := range matches {
			match.print(os.Stdout, MatchOutputHighlight)
		}
	}
	return 0
}

// GetSearchParams returns the search parameters the command was run with
func (cmd *GrepCommand) GetSearchParams() SearchParameters {
	return SearchParameters{
		Search:   cmd.args.Search,
		Mode:     cmd.Mode,
		IsRegexp: cmd.args.Regexp,
	}
}

func (cmd *GrepCommand) grepFile(search string, path string) (matches []*Match, err error) {
	matches = []*Match{}

	if cmd.client.GetType(path) == client.LEAF {
		secret, err := cmd.client.Read(path)
		if err != nil {
			return matches, err
		}

		for k, v := range secret.GetData() {
			matches = append(matches, cmd.searcher.DoSearch(path, k, fmt.Sprintf("%v", v))...)
		}
	}

	return matches, nil
}
