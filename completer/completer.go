package completer

import (
	"github.com/fishi0x01/vsh/log"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/structs"
	"github.com/fishi0x01/vsh/cli"
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

func (c *Completer) isCommandArgument(p string) bool {
	words := strings.Split(p, " ")
	if len(words) < 2 {
		return false
	}

	commands := cli.NewCommands(c.client)
	for _, f := range structs.Fields(commands) {
		if words[0] == f.Value().(cli.Command).GetName() || words[0] == "toggle-auto-completion" {
			return true
		}
	}
	return false
}

func isCommand(p string) bool {
	return len(strings.Split(p, " ")) < 2
}

func (c *Completer) commandSuggestions(arg string) (result []prompt.Suggest) {
	result = make([]prompt.Suggest, 0)
	commands := cli.NewCommands(c.client)
	for _, f := range structs.Fields(commands) {
		val := f.Value().(cli.Command)
		result = append(result, prompt.Suggest{Text: val.GetName(), Description: cli.Usage(val)})
	}
	result = append(result, prompt.Suggest{Text: "toggle-auto-completion", Description: "toggle path auto-completion on/off"})

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
	} else if c.isCommandArgument(p) && c.pathCompletionToggle {
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
