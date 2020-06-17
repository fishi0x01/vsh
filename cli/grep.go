package cli

import (
	"fmt"
	"index/suffixarray"
	"io"
	"sort"

	"github.com/fatih/color"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// GrepCommand container for all 'grep' parameters
type GrepCommand struct {
	name string

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Path   string
	Search string
}

// Match structure to keep indices of matched terms
type Match struct {
	path  string
	term  string
	key   string
	value string
	// sorted slices of indices of match starts
	keyIndex   []int
	valueIndex []int
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

// Parse given arguments and return status
func (cmd *GrepCommand) Parse(args []string) (success bool) {
	if len(args) == 3 {
		cmd.Search = args[1]
		cmd.Path = args[2]
		success = true
	} else {
		fmt.Println("Usage:\ngrep <term-string> <path>")
	}
	return success
}

// Run executes 'grep' with given RemoveCommand's parameters
func (cmd *GrepCommand) Run() {
	path := cmdPath(cmd.client.Pwd, cmd.Path)

	t := cmd.client.GetType(path)
	if t != client.NODE && t != client.LEAF {
		log.Error("Invalid path: %s", path)
		return
	}

	for _, path := range cmd.client.Traverse(path) {
		matches, err := cmd.grepFile(cmd.Search, path)
		if err != nil {
			return
		}
		for _, match := range matches {
			match.print(cmd.stdout)
		}
	}

	return
}

func (cmd *GrepCommand) grepFile(search string, path string) (matches []*Match, err error) {
	matches = []*Match{}

	if cmd.client.GetType(path) == client.LEAF {
		secret, err := cmd.client.Read(path)
		if err != nil {
			return matches, err
		}

		for k, v := range secret.Data {
			if rec, ok := v.(map[string]interface{}); ok {
				// KV 2
				for kk, vv := range rec {
					matches = append(matches, match(path, kk, fmt.Sprintf("%v", vv), search)...)
				}
			} else {
				// KV 1
				matches = append(matches, match(path, k, fmt.Sprintf("%v", v), search)...)
			}
		}
	}

	return matches, nil
}

// find all indices for matches in key and value
func match(path string, k string, v string, substr string) (m []*Match) {
	keyIndex := suffixarray.New([]byte(k))
	keyMatches := keyIndex.Lookup([]byte(substr), -1)
	sort.Ints(keyMatches)

	valueIndex := suffixarray.New([]byte(v))
	valueMatches := valueIndex.Lookup([]byte(substr), -1)
	sort.Ints(valueMatches)

	if len(keyMatches) > 0 || len(valueMatches) > 0 {
		m = []*Match{
			&Match{
				path:       path,
				term:       substr,
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
	highlightMatches(match.key, match.term, match.keyIndex, out)
	fmt.Fprintf(out, " : ")
	highlightMatches(match.value, match.term, match.valueIndex, out)
	fmt.Fprintf(out, "\n")
}

func highlightMatches(s string, term string, index []int, out io.Writer) {
	matchColor := color.New(color.FgYellow).SprintFunc()
	cur := 0
	if len(index) > 0 {
		for _, next := range index {
			end := next + len(term)
			fmt.Fprint(out, s[cur:next])
			fmt.Fprint(out, matchColor(s[next:end]))
			cur = end
		}
		fmt.Fprint(out, s[cur:])
	} else {
		fmt.Fprint(out, s)
	}
}
