package completer

import (
	"github.com/fishi0x01/vsh/log"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fishi0x01/vsh/client"
)

// Completer struct for tab completion
type Completer struct {
	pathCompletionToggle bool
	client               *client.Client
}

// NewCompleter creates a new Completer with given client
func NewCompleter(client *client.Client, disableAutoCompletion bool) *Completer {
	return &Completer{
		pathCompletionToggle: !disableAutoCompletion,
		client:               client,
	}
}

// TogglePathCompletion enable/disable path auto-completion
func (c *Completer) TogglePathCompletion() {
	c.pathCompletionToggle = !c.pathCompletionToggle
	log.UserInfo("Use path auto-completion: %t", c.pathCompletionToggle)
}

func (c *Completer) getAbsoluteTopLevelSuggestions() []prompt.Suggest {
	var suggestions []prompt.Suggest
	for k := range c.client.KVBackends {
		suggestions = append(suggestions, prompt.Suggest{Text: "/" + k})
	}
	return suggestions
}

func (c *Completer) getRelativeTopLevelSuggestions() []prompt.Suggest {
	var suggestions []prompt.Suggest
	for k := range c.client.KVBackends {
		suggestions = append(suggestions, prompt.Suggest{Text: k})
	}
	return suggestions
}

func (c *Completer) absolutePathSuggestions(arg string) (result []prompt.Suggest) {
	if strings.Count(arg, "/") < 2 {
		result = c.getAbsoluteTopLevelSuggestions()
	} else {
		li := strings.LastIndex(arg, "/")
		queryPath := arg[0 : li+1]

		var options []string
		var err error
		options, err = c.client.List(queryPath)

		if err != nil {
			panic(err)
		}

		options = append(options, "../")
		for _, node := range options {
			result = append(result, prompt.Suggest{Text: queryPath + node})
		}
	}

	filtered := prompt.FilterHasPrefix(result, arg, true)
	if len(filtered) > 0 {
		result = filtered
	}
	return result
}

func (c *Completer) relativePathSuggestions(arg string) (result []prompt.Suggest) {
	if c.client.Pwd == "/" && strings.Count(arg, "/") < 1 {
		result = c.getRelativeTopLevelSuggestions()
	} else {
		li := strings.LastIndex(arg, "/")
		queryPath := arg[0 : li+1]

		var options []string
		var err error
		options, err = c.client.List(c.client.Pwd + queryPath)

		if err != nil {
			return result
		}

		options = append(options, "../")
		for _, node := range options {
			result = append(result, prompt.Suggest{Text: queryPath + node})
		}
	}

	filtered := prompt.FilterHasPrefix(result, arg, true)
	if len(filtered) > 0 {
		result = filtered
	}
	return result
}

func isAbsolutePath(path string) bool {
	return strings.HasPrefix(path, "/")
}

func isCommandArgument(p string) bool {
	words := strings.Split(p, " ")
	if len(words) < 2 {
		return false
	}

	return words[0] == "cd" ||
		words[0] == "cp" ||
		words[0] == "rm" ||
		words[0] == "mv" ||
		words[0] == "grep" ||
		words[0] == "cat" ||
		words[0] == "append" ||
		words[0] == "ls" ||
		words[0] == "toggle-auto-completion"
}

func isCommand(p string) bool {
	return len(strings.Split(p, " ")) < 2
}

func (c *Completer) commandSuggestions(arg string) (result []prompt.Suggest) {
	result = []prompt.Suggest{
		{Text: "append", Description: "append <from> <to> [-f|--force] | [-s|--skip] | [-r|--rename] | -s is default"},
		{Text: "cat", Description: "cat <path>"},
		{Text: "cd", Description: "cd <path>"},
		{Text: "cp", Description: "cp <from> <to> | -r is implied"},
		{Text: "grep", Description: "grep <search> <path> [-e|--regexp] [-k|--keys] [-v|--values]"},
		{Text: "ls", Description: "ls <path>"},
		{Text: "mv", Description: "mv <from> <to>"},
		{Text: "replace", Description: "replace <search> <replacement> <path> [-e|--regexp] [-k|--keys] [-v|--values] [-y|--confirm] [-n|--dry-run]"},
		{Text: "rm", Description: "rm <path> | -r is implied"},
		{Text: "toggle-auto-completion", Description: "toggle path auto-completion on/off"},
	}
	filtered := prompt.FilterHasPrefix(result, arg, true)
	if len(filtered) > 0 {
		result = filtered
	}
	return result
}

// Complete suggestions for completion
func (c *Completer) Complete(in prompt.Document) (result []prompt.Suggest) {
	p := in.TextBeforeCursor()
	if isCommand(p) {
		result = c.commandSuggestions(in.GetWordBeforeCursor())
	} else if isCommandArgument(p) && c.pathCompletionToggle {
		cur := in.GetWordBeforeCursor()
		if isAbsolutePath(cur) {
			result = c.absolutePathSuggestions(cur)
		} else {
			result = c.relativePathSuggestions(cur)
		}
	}

	return result
}

// PromptPrefix returns the currently active prompt prefix
func (c *Completer) PromptPrefix() (string, bool) {
	return c.client.Name + " " + c.client.Pwd + "> ", true
}
