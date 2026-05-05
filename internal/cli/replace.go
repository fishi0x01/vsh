package cli

import (
	"fmt"
	"os"
	"sync"

	"github.com/fishi0x01/vsh/internal/client"
	"github.com/fishi0x01/vsh/internal/logger"
)

// ReplaceCommand container for all 'replace' parameters
type ReplaceCommand struct {
	name        string
	args        *ReplaceCommandArgs
	workerCount int

	client        *client.Client
	searcher      *Searcher
	Mode          KeyValueMode
	pauseSpinner  func()
	resumeSpinner func()
}

// SetSpinnerControl satisfies SpinnerAware so the spinner can be paused
// while askForConfirmation is reading stdin.
func (cmd *ReplaceCommand) SetSpinnerControl(pause, resume func()) {
	cmd.pauseSpinner = pause
	cmd.resumeSpinner = resume
}

// ReplaceCommandArgs provides a struct for go-arg parsing
type ReplaceCommandArgs struct {
	Search      string         `arg:"positional,required"`
	Replacement string         `arg:"positional,required"`
	Path        string         `arg:"positional"`
	Confirm     bool           `arg:"-y,--confirm"        help:"Write results without prompt"`
	DryRun      bool           `arg:"-n,--dry-run"        help:"Skip writing results without prompt"`
	KeySelector string         `arg:"-s,--key-selector"   help:"Limit replacements to specified key"                           placeholder:"PATTERN"`
	Keys        bool           `arg:"-k,--keys"           help:"Match against keys (true if -v is not specified)"`
	Output      MatchOutputArg `arg:"-o,--output"         help:"Present changes as 'inline' with color or traditional 'diff'"                        default:"inline"`
	Regexp      bool           `arg:"-e,--regexp"         help:"Treat search string and selector as a regexp"`
	Shallow     bool           `arg:"-S,--shallow"        help:"Only search leaf nodes of the path rather than recurse deeper"`
	Values      bool           `arg:"-v,--values"         help:"Match against values (true if -k is not specified)"`
}

// Description provides detail on what the command does
func (ReplaceCommandArgs) Description() string {
	return "recursively replaces a matching pattern with a replacement string at a path"
}

// NewReplaceCommand creates a new ReplaceCommand parameter container
func NewReplaceCommand(c *client.Client, workerCount int) *ReplaceCommand {
	return &ReplaceCommand{
		name:        "replace",
		client:      c,
		args:        &ReplaceCommandArgs{},
		workerCount: workerCount,
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
	if cmd.args.Keys {
		cmd.Mode |= ModeKeys
	}
	if cmd.args.Values {
		cmd.Mode |= ModeValues
	}
	if cmd.Mode == 0 {
		cmd.Mode = ModeKeys + ModeValues
	}
	if cmd.args.DryRun && cmd.args.Confirm {
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
	filePaths, err := cmd.client.SubpathsForPath(path, cmd.args.Shallow)
	if err != nil {
		logger.UserError(fmt.Sprintf("%s", err))
		return 1
	}

	allMatches, err := cmd.findMatches(filePaths)
	if err != nil {
		logger.UserError(fmt.Sprintf("%s", err))
		return 1
	}
	return cmd.commitMatches(allMatches)
}

func (cmd *ReplaceCommand) findMatches(
	filePaths []string,
) (matchesByPath map[string][]*Match, err error) {
	type result struct {
		path    string
		matches []*Match
		err     error
	}
	results := make([]result, len(filePaths))

	var wg sync.WaitGroup
	sem := make(chan struct{}, cmd.workerCount)
	for i, curPath := range filePaths {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			matches, err := cmd.FindReplacements(cmd.args.Search, cmd.args.Replacement, p)
			results[idx] = result{p, matches, err}
		}(i, curPath)
	}
	wg.Wait()

	matchesByPath = make(map[string][]*Match)
	for _, r := range results {
		if r.err != nil {
			return matchesByPath, r.err
		}
		for _, match := range r.matches {
			match.print(os.Stdout, cmd.args.Output.Value)
		}
		if len(r.matches) > 0 {
			matchesByPath[r.path] = append(matchesByPath[r.path], r.matches...)
		}
	}
	return matchesByPath, nil
}

func (cmd *ReplaceCommand) commitMatches(matchesByPath map[string][]*Match) int {
	if len(matchesByPath) > 0 {
		if !cmd.args.Confirm && !cmd.args.DryRun {
			if cmd.pauseSpinner != nil {
				cmd.pauseSpinner()
			}
			result, err := askForConfirmation("Write changes to Vault?")
			if cmd.resumeSpinner != nil {
				cmd.resumeSpinner()
			}
			if err != nil {
				return 1
			}
			cmd.args.Confirm = result
		}
		if !cmd.args.Confirm {
			fmt.Println("Skipping write.")
			return 0
		}
		fmt.Println("Writing!")
		err := cmd.WriteReplacements(matchesByPath)
		if err != nil {
			logger.UserError("Error writing replacement: %v", err)
		}
	} else {
		fmt.Println("No matches found to replace.")
	}
	return 0
}

// FindReplacements will find the matches for a given search string to be replaced
func (cmd *ReplaceCommand) FindReplacements(
	search string,
	replacement string,
	path string,
) (matches []*Match, err error) {
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

// WriteReplacements will write replacement data back to Vault concurrently
func (cmd *ReplaceCommand) WriteReplacements(groupedMatches map[string][]*Match) error {
	failed := make(chan error, 1)
	var wg sync.WaitGroup
	sem := make(chan struct{}, cmd.workerCount)

	for path, matches := range groupedMatches {
		wg.Add(1)
		go func(p string, m []*Match) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			secret, err := cmd.client.Read(p)
			if err != nil {
				select {
				case failed <- err:
				default:
				}
				return
			}
			data := secret.GetData()
			for _, match := range m {
				if p != match.path {
					select {
					case failed <- fmt.Errorf("match path does not equal group path"):
					default:
					}
					return
				}
				if match.replacedKey != match.key {
					delete(data, match.key)
				}
				data[match.replacedKey] = match.replacedValue
			}
			secret.SetData(data)
			if err := cmd.client.Write(p, secret); err != nil {
				select {
				case failed <- err:
				default:
				}
			}
		}(path, matches)
	}
	wg.Wait()

	select {
	case err := <-failed:
		return err
	default:
		return nil
	}
}
