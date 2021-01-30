package cli

import (
	"fmt"
	"index/suffixarray"
	"io"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/fatih/color"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// GrepMode defines the scope of which parts of a path to search (keys and/or values)
type GrepMode int

const (
	// ModeKeys only searches keys
	ModeKeys GrepMode = 1
	// ModeValues only searches values
	ModeValues GrepMode = 2
)

// GrepCommand container for all 'grep' parameters
type GrepCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Path   string
	Search string
	Regexp *regexp.Regexp
	Mode   GrepMode
}

// Match structure to keep indices of matched terms
type Match struct {
	path  string
	term  string
	key   string
	value string
	// sorted slices of indices of match starts and length
	keyIndex   [][]int
	valueIndex [][]int
}

// NewGrepCommand creates a new GrepCommand parameter container
func NewGrepCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *GrepCommand {
	return &GrepCommand{
		name:   "grep",
		client: c,
		stdout: stdout,
		stderr: stderr,
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

// IsMode returns true if the specified mode is enabled
func (cmd *GrepCommand) IsMode(mode GrepMode) bool {
	return cmd.Mode&mode == mode
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
			re, err := regexp.Compile(cmd.Search)
			if err != nil {
				return fmt.Errorf("cannot parse regex pattern")
			}
			cmd.Regexp = re
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
			match.print(cmd.stdout)
		}
	}
	return 0
}

func (cmd *GrepCommand) grepFile(search string, path string) (matches []*Match, err error) {
	matches = []*Match{}

	if cmd.client.GetType(path) == client.LEAF {
		secret, err := cmd.client.Read(path)
		if err != nil {
			return matches, err
		}

		for k, v := range secret.GetData() {
			matches = append(matches, cmd.doMatch(path, k, fmt.Sprintf("%v", v), search)...)
		}
	}

	return matches, nil
}

func (cmd *GrepCommand) doMatch(path string, k string, v string, search string) (m []*Match) {
	if cmd.Regexp != nil {
		return cmd.regexpMatch(path, k, v, cmd.Regexp)
	}
	return cmd.substrMatch(path, k, v, search)
}

// find all indices for matches in key and value
func (cmd *GrepCommand) substrMatch(path string, k string, v string, substr string) (m []*Match) {
	substrLength := len(substr)
	keyMatchPairs := make([][]int, 0)
	if cmd.IsMode(ModeKeys) {
		keyIndex := suffixarray.New([]byte(k))
		keyMatches := keyIndex.Lookup([]byte(substr), -1)
		sort.Ints(keyMatches)
		for _, offset := range keyMatches {
			keyMatchPairs = append(keyMatchPairs, []int{offset, substrLength})
		}
	}

	valueMatchPairs := make([][]int, 0)
	if cmd.IsMode(ModeValues) {
		valueIndex := suffixarray.New([]byte(v))
		valueMatches := valueIndex.Lookup([]byte(substr), -1)
		sort.Ints(valueMatches)
		for _, offset := range valueMatches {
			valueMatchPairs = append(valueMatchPairs, []int{offset, substrLength})
		}
	}

	if len(keyMatchPairs) > 0 || len(valueMatchPairs) > 0 {
		m = []*Match{
			{
				path:       path,
				term:       substr,
				key:        k,
				value:      v,
				keyIndex:   keyMatchPairs,
				valueIndex: valueMatchPairs,
			},
		}
	}
	return m
}

func (cmd *GrepCommand) regexpMatch(path string, k string, v string, pattern *regexp.Regexp) (m []*Match) {
	keyMatches := make([][]int, 0)
	if cmd.IsMode(ModeKeys) {
		keyMatches = pattern.FindAllIndex([]byte(k), -1)
	}
	valueMatches := make([][]int, 0)
	if cmd.IsMode(ModeValues) {
		valueMatches = pattern.FindAllIndex([]byte(v), -1)
	}

	if len(keyMatches) > 0 || len(valueMatches) > 0 {
		m = []*Match{
			{
				path:       path,
				term:       pattern.String(),
				key:        k,
				value:      v,
				keyIndex:   keyMatches,
				valueIndex: valueMatches,
			},
		}
	}
	return m
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
