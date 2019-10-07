package completer

import (
	"github.com/c-bata/go-prompt"
	"github.com/fishi0x01/vsh/client"
	"strings"
)

// Completer struct for tab completion
type Completer struct {
	client *client.Client
}

// NewCompleter creates a new Completer with given client
func NewCompleter(client *client.Client) *Completer {
	return &Completer{
		client: client,
	}
}

func (c *Completer) getAbsoluteTopLevelSuggestions() []prompt.Suggest {
	var suggestions = []prompt.Suggest{}
	for k := range c.client.KVBackends {
		suggestions = append(suggestions, prompt.Suggest{"/" + k, ""})
	}
	return suggestions
}

func (c *Completer) getRelativeTopLevelSuggestions() []prompt.Suggest {
	var suggestions = []prompt.Suggest{}
	for k := range c.client.KVBackends {
		suggestions = append(suggestions, prompt.Suggest{k, ""})
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
			result = append(result, prompt.Suggest{queryPath + node, ""})
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
			panic(err)
		}

		options = append(options, "../")
		for _, node := range options {
			result = append(result, prompt.Suggest{queryPath + node, ""})
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
		words[0] == "cat" ||
		words[0] == "ls"
}

func isCommand(p string) bool {
	return len(strings.Split(p, " ")) < 2
}

func (c *Completer) commandSuggestions(arg string) (result []prompt.Suggest) {
	result = []prompt.Suggest{
		{Text: "cd", Description: "Change directory"},
		{Text: "cp", Description: "Copy secret or directory"},
		{Text: "rm", Description: "Remove secret or directory"},
		{Text: "mv", Description: "Move secret or directory"},
		{Text: "cat", Description: "Print secret content"},
		{Text: "ls", Description: "List secrets and directories"},
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
	} else if isCommandArgument(p) {
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
