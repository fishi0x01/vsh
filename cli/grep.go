package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// GrepCommand container for all 'grep' parameters
type GrepCommand struct {
	name string

	client   *client.Client
	Path     string
	searcher *Searcher
	SearchParameters
}

// NewGrepCommand creates a new GrepCommand parameter container
func NewGrepCommand(c *client.Client) *GrepCommand {
	return &GrepCommand{
		name:   "grep",
		client: c,
	}
}

// GetName returns the GrepCommand's name identifier
func (cmd *GrepCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *GrepCommand) IsSane() bool {
	return cmd.Path != "" && cmd.Search != ""
}

// PrintUsage print command usage
func (cmd *GrepCommand) PrintUsage() {
	log.UserInfo("Usage:\ngrep <search> <path> [-e|--regexp] [-k|--keys] [-v|--values]")
}

// Parse given arguments and return status
func (cmd *GrepCommand) Parse(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Search = args[1]
	cmd.Path = args[2]
	flags := args[3:]

	for _, v := range flags {
		switch v {
		case "-e", "--regexp":
			cmd.doRegexp = true
		case "-k", "--keys":
			cmd.Mode |= ModeKeys
		case "-v", "--values":
			cmd.Mode |= ModeValues
		default:
			return fmt.Errorf("invalid flag: %s", v)
		}
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

// Run executes 'grep' with given RemoveCommand's parameters
func (cmd *GrepCommand) Run() int {
	path := cmdPath(cmd.client.Pwd, cmd.Path)
	var filePaths []string

	switch t := cmd.client.GetType(path); t {
	case client.LEAF:
		filePaths = append(filePaths, filepath.Clean(path))
	case client.NODE:
		for _, traversedPath := range cmd.client.Traverse(path) {
			filePaths = append(filePaths, traversedPath)
		}
	default:
		log.UserError("Not a valid path for operation: %s", path)
		return 1
	}

	for _, curPath := range filePaths {
		matches, err := cmd.grepFile(cmd.Search, curPath)
		if err != nil {
			return 1
		}
		for _, match := range matches {
			match.print(os.Stdout)
		}
	}
	return 0
}

func (cmd *GrepCommand) GetSearchParams() SearchParameters {
	return SearchParameters{
		Search: cmd.Search,
		Mode: cmd.Mode,
		doRegexp: cmd.doRegexp,
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
			matches = append(matches, cmd.searcher.DoMatch(path, k, fmt.Sprintf("%v", v))...)
		}
	}

	return matches, nil
}

func (match *Match) print(out io.Writer) {
	fmt.Fprint(out, match.path, "> ")
	highlightMatches(match.key, match.keyIndex, out)
	fmt.Fprintf(out, " = ")
	highlightMatches(match.value, match.valueIndex, out)
	fmt.Fprintf(out, "\n")
}

func highlightMatches(s string, index [][]int, out io.Writer) {
	matchColor := color.New(color.FgYellow).SprintFunc()
	cur := 0
	if len(index) > 0 {
		for _, pair := range index {
			next := pair[0]
			length := pair[1]
			end := next + length
			fmt.Fprint(out, s[cur:next])
			fmt.Fprint(out, matchColor(s[next:end]))
			cur = end
		}
		fmt.Fprint(out, s[cur:])
	} else {
		fmt.Fprint(out, s)
	}
}
