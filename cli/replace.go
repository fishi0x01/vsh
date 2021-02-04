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

	client   *client.Client
	Confirm  bool
	DryRun   bool
	Path     string
	searcher *Searcher
	SearchParameters
}

// NewReplaceCommand creates a new ReplaceCommand parameter container
func NewReplaceCommand(c *client.Client) *ReplaceCommand {
	return &ReplaceCommand{
		name:   "replace",
		client: c,
	}
}

// GetName returns the ReplaceCommand's name identifier
func (cmd *ReplaceCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *ReplaceCommand) IsSane() bool {
	return cmd.Search != "" && cmd.Replacement != nil && cmd.Path != ""
}

// PrintUsage print command usage
func (cmd *ReplaceCommand) PrintUsage() {
	log.UserInfo("Usage:\nreplace <search> <replacement> <path> [-e|--regexp] [-k|--keys] [-v|--values]")
}

// GetSearchParams returns the search parameters the command was run with
func (cmd *ReplaceCommand) GetSearchParams() SearchParameters {
	return SearchParameters{
		Search:      cmd.Search,
		Replacement: cmd.Replacement,
		Mode:        cmd.Mode,
		IsRegexp:    cmd.IsRegexp,
	}
}

// Parse given arguments and return status
func (cmd *ReplaceCommand) Parse(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Search = args[1]
	cmd.Replacement = &args[2]
	cmd.Path = args[3]
	flags := args[4:]

	for _, v := range flags {
		switch v {
		case "-e", "--regexp":
			cmd.IsRegexp = true
		case "-k", "--keys":
			cmd.Mode |= ModeKeys
		case "-v", "--values":
			cmd.Mode |= ModeValues
		case "-n", "--dry-run":
			cmd.DryRun = true
		case "-y", "--confirm":
			cmd.Confirm = true
		default:
			return fmt.Errorf("invalid flag: %s", v)
		}
	}
	if cmd.Mode == 0 {
		cmd.Mode = ModeKeys + ModeValues
	}
	if cmd.DryRun == true && cmd.Confirm == true {
		cmd.Confirm = false
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
	path := cmdPath(cmd.client.Pwd, cmd.Path)
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
		matches, err := cmd.FindReplacements(cmd.Search, *cmd.Replacement, curPath)
		if err != nil {
			return matchesByPath, err
		}
		for _, match := range matches {
			match.print(os.Stdout, true)
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
		if cmd.Confirm == false && cmd.DryRun == false {
			p := promptx.NewDefaultConfirm("Write changes to Vault?", false)
			result, err := p.Run()
			if err != nil {
				return 1
			}
			cmd.Confirm = result
		}
		if cmd.Confirm == false {
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
