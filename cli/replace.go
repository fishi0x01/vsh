package cli

import (
	"fmt"
	"os"

	"github.com/cnlubo/promptx"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
)

// ReplaceCommand container for all 'replace' parameters
type ReplaceCommand struct {
	name string
	args *ReplaceCommandArgs

	client   *client.Client
	searcher *Searcher
	Mode     KeyValueMode
}

// ReplaceCommandArgs provides a struct for go-arg parsing
type ReplaceCommandArgs struct {
	Search      string         `arg:"positional,required"`
	Replacement string         `arg:"positional,required"`
	Path        string         `arg:"positional"`
	Regexp      bool           `arg:"-e,--regexp" help:"Treat search string and selector as a regexp"`
	KeySelector string         `arg:"-s,--key-selector" help:"Limit replacements to specified key" placeholder:"PATTERN"`
	Keys        bool           `arg:"-k,--keys" help:"Match against keys (true if -v is not specified)"`
	Values      bool           `arg:"-v,--values" help:"Match against values (true if -k is not specified)"`
	Confirm     bool           `arg:"-y,--confirm" help:"Write results without prompt"`
	DryRun      bool           `arg:"-n,--dry-run" help:"Skip writing results without prompt"`
	Output      MatchOutputArg `arg:"-o,--output" help:"Present changes as 'inline' with color or traditional 'diff'" default:"inline"`
}

// Description provides detail on what the command does
func (ReplaceCommandArgs) Description() string {
	return "recursively replaces a matching pattern with a replacement string at a path"
}

// NewReplaceCommand creates a new ReplaceCommand parameter container
func NewReplaceCommand(c *client.Client) *ReplaceCommand {
	return &ReplaceCommand{
		name:   "replace",
		client: c,
		args:   &ReplaceCommandArgs{},
	}
}

// GetName returns the ReplaceCommand's name identifier
func (cmd *ReplaceCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *ReplaceCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *ReplaceCommand) IsSane() bool {
	return cmd.args.Search != "" && cmd.args.Path != ""
}

// PrintUsage print command usage
func (cmd *ReplaceCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// GetSearchParams returns the search parameters the command was run with
func (cmd *ReplaceCommand) GetSearchParams() SearchParameters {
	return SearchParameters{
		IsRegexp:    cmd.args.Regexp,
		KeySelector: cmd.args.KeySelector,
		Mode:        cmd.Mode,
		Output:      cmd.args.Output.Value,
		Replacement: &cmd.args.Replacement,
		Search:      cmd.args.Search,
	}
}

// Parse given arguments and return status
func (cmd *ReplaceCommand) Parse(args []string) error {
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
	if cmd.args.DryRun == true && cmd.args.Confirm == true {
		cmd.args.Confirm = false
	}

	searcher, err := NewSearcher(cmd)
	if err != nil {
		return err
	}
	cmd.searcher = searcher

	return nil
}

// Run executes 'replace' with given ReplaceCommand's parameters
func (cmd *ReplaceCommand) Run() int {
	path := cmdPath(cmd.client.Pwd, cmd.args.Path)
	filePaths, err := cmd.client.SubpathsForPath(path)
	if err != nil {
		log.UserError(fmt.Sprintf("%s", err))
		return 1
	}

	allMatches, err := cmd.findMatches(filePaths)
	if err != nil {
		log.UserError(fmt.Sprintf("%s", err))
		return 1
	}
	return cmd.commitMatches(allMatches)
}

func (cmd *ReplaceCommand) findMatches(filePaths []string) (matchesByPath map[string][]*Match, err error) {
	matchesByPath = make(map[string][]*Match, 0)
	for _, curPath := range filePaths {
		matches, err := cmd.FindReplacements(cmd.args.Search, cmd.args.Replacement, curPath)
		if err != nil {
			return matchesByPath, err
		}
		for _, match := range matches {
			match.print(os.Stdout, cmd.args.Output.Value)
		}
		if len(matches) > 0 {
			_, ok := matchesByPath[curPath]
			if ok == false {
				matchesByPath[curPath] = make([]*Match, 0)
			}
			matchesByPath[curPath] = append(matchesByPath[curPath], matches...)
		}
	}
	return matchesByPath, nil
}

func (cmd *ReplaceCommand) commitMatches(matchesByPath map[string][]*Match) int {
	if len(matchesByPath) > 0 {
		if cmd.args.Confirm == false && cmd.args.DryRun == false {
			p := promptx.NewDefaultConfirm("Write changes to Vault?", false)
			result, err := p.Run()
			if err != nil {
				return 1
			}
			cmd.args.Confirm = result
		}
		if cmd.args.Confirm == false {
			fmt.Println("Skipping write.")
			return 0
		}
		fmt.Println("Writing!")
		cmd.WriteReplacements(matchesByPath)
	} else {
		fmt.Println("No matches found to replace.")
	}
	return 0
}

// FindReplacements will find the matches for a given search string to be replaced
func (cmd *ReplaceCommand) FindReplacements(search string, replacement string, path string) (matches []*Match, err error) {
	if cmd.client.GetType(path) == client.LEAF {
		secret, err := cmd.client.Read(path)
		if err != nil {
			return matches, err
		}

		for k, v := range secret.GetData() {
			match := cmd.searcher.DoSearch(path, k, fmt.Sprintf("%v", v))
			matches = append(matches, match...)
		}
	}
	return matches, nil
}

// WriteReplacements will write replacement data back to Vault
func (cmd *ReplaceCommand) WriteReplacements(groupedMatches map[string][]*Match) error {
	// process matches by vault path
	for path, matches := range groupedMatches {
		secret, err := cmd.client.Read(path)
		if err != nil {
			return err
		}
		data := secret.GetData()

		// update secret with changes. remove key w/ prior names, add renamed keys, update values.
		for _, match := range matches {
			if path != match.path {
				return fmt.Errorf("match path does not equal group path")
			}
			if match.replacedKey != match.key {
				delete(data, match.key)
			}
			data[match.replacedKey] = match.replacedValue
		}
		secret.SetData(data)

		err = cmd.client.Write(path, secret)
		if err != nil {
			return err
		}
	}
	return nil
}
