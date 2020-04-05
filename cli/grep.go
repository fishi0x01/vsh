package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fishi0x01/vsh/client"
	"index/suffixarray"
	"io"
	"sort"
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

type Match struct {
	path   string
	search string
	key    string
	value  string
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
		fmt.Println("Usage:\ngrep <search-string> <path>")
	}
	return success
}

// Run executes 'grep' with given RemoveCommand's parameters
func (cmd *GrepCommand) Run() {
	path := cmdPath(cmd.client.Pwd, cmd.Path)

	t := cmd.client.GetType(path)
	if t != client.NODE && t != client.LEAF {
		fmt.Fprintln(cmd.stderr, "Not a valid path: "+path)
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
					m := match(path, kk, fmt.Sprintf("%v", vv), search)
					if m != nil {
						matches = append(matches, m)
					}
				}
			} else {
				// KV 1
				m := match(path, k, fmt.Sprintf("%v", v), search)
				if m != nil {
					matches = append(matches, m)
				}
			}
		}
	}

	return matches, nil
}

// find all indices for matches in key and value
func match(path string, k string, v string, substr string) (m *Match) {
	keyIndex := suffixarray.New([]byte(k))
	keyMatches := keyIndex.Lookup([]byte(substr), -1)
	sort.Ints(keyMatches)

	valueIndex := suffixarray.New([]byte(v))
	valueMatches := valueIndex.Lookup([]byte(substr), -1)
	sort.Ints(valueMatches)

	if len(keyMatches) > 0 || len(valueMatches) > 0 {
		m = &Match{
			path:       path,
			search:     substr,
			key:        k,
			value:      v,
			keyIndex:   keyMatches,
			valueIndex: valueMatches,
		}
	}

	return m
}

func (match *Match) print(out io.Writer) {
	matchColor := color.New(color.FgYellow).SprintFunc()
	fmt.Fprint(out, match.path, "> ")

	cur := 0
	if len(match.keyIndex) > 0 {
		for _, index := range match.keyIndex {
			end := index + len(match.search)
			fmt.Fprint(out, match.key[cur:index])
			fmt.Fprint(out, matchColor(match.key[index:end]))
			cur = end
		}
		fmt.Fprint(out, match.key[cur:], " : ")
	} else {
		fmt.Fprint(out, match.key, " : ")
	}

	cur = 0
	if len(match.valueIndex) > 0 {
		for _, index := range match.valueIndex {
			end := index + len(match.search)
			fmt.Fprint(out, match.value[cur:index])
			fmt.Fprint(out, matchColor(match.value[index:end]))
			cur = end
		}
		fmt.Fprint(out, match.value[cur:], "\n")
	} else {
		fmt.Fprint(out, match.value, "\n")
	}
}
