package completer

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/fishi0x01/vsh/internal/cli"
	"github.com/fishi0x01/vsh/internal/client"
	"github.com/fishi0x01/vsh/internal/logger"
)

// Suggestion holds a completion candidate with an optional description.
type Suggestion struct {
	Text        string
	Description string
}

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
	logger.UserInfo("Use path auto-completion: %t", c.pathCompletionToggle)
}

func filterHasPrefix(suggestions []Suggestion, prefix string) []Suggestion {
	if prefix == "" {
		return suggestions
	}
	lower := strings.ToLower(prefix)
	var result []Suggestion
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(s.Text), lower) {
			result = append(result, s)
		}
	}
	return result
}

func (c *Completer) getAbsoluteTopLevelSuggestions() []Suggestion {
	var suggestions []Suggestion
	for k := range c.client.KVBackends {
		suggestions = append(suggestions, Suggestion{Text: "/" + k})
	}
	return suggestions
}

func (c *Completer) getRelativeTopLevelSuggestions() []Suggestion {
	var suggestions []Suggestion
	for k := range c.client.KVBackends {
		suggestions = append(suggestions, Suggestion{Text: k})
	}
	return suggestions
}

func (c *Completer) absolutePathSuggestions(arg string) (result []Suggestion) {
	if strings.Count(arg, "/") < 2 {
		result = c.getAbsoluteTopLevelSuggestions()
	} else {
		li := strings.LastIndex(arg, "/")
		queryPath := arg[0 : li+1]

		options, err := c.client.List(queryPath)
		if err != nil {
			return result
		}

		options = append(options, "../")
		for _, node := range options {
			result = append(result, Suggestion{Text: queryPath + node})
		}
	}

	filtered := filterHasPrefix(result, arg)
	if len(filtered) > 0 {
		result = filtered
	}
	return result
}

func (c *Completer) relativePathSuggestions(arg string) (result []Suggestion) {
	if c.client.Pwd == "/" && strings.Count(arg, "/") < 1 {
		result = c.getRelativeTopLevelSuggestions()
	} else {
		li := strings.LastIndex(arg, "/")
		queryPath := arg[0 : li+1]

		options, err := c.client.List(c.client.Pwd + queryPath)
		if err != nil {
			return result
		}

		options = append(options, "../")
		for _, node := range options {
			result = append(result, Suggestion{Text: queryPath + node})
		}
	}

	filtered := filterHasPrefix(result, arg)
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

	commands := cli.NewCommands(c.client, 1, false)
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

func (c *Completer) commandSuggestions(arg string) []Suggestion {
	result := make([]Suggestion, 0)
	commands := cli.NewCommands(c.client, 1, false)
	for _, f := range structs.Fields(commands) {
		val := f.Value().(cli.Command)
		result = append(result, Suggestion{Text: val.GetName(), Description: cli.Usage(val)})
	}
	result = append(result, Suggestion{
		Text:        "toggle-auto-completion",
		Description: "toggle path auto-completion on/off",
	})

	filtered := filterHasPrefix(result, arg)
	if len(filtered) > 0 {
		result = filtered
	}
	return result
}

// Complete returns completion suggestions for the given input text.
func (c *Completer) Complete(text string) []Suggestion {
	if isCommand(text) {
		// word before cursor is the partial command name
		word := text
		if idx := strings.LastIndex(text, " "); idx >= 0 {
			word = text[idx+1:]
		}
		return c.commandSuggestions(word)
	}

	if c.isCommandArgument(text) && c.pathCompletionToggle {
		// word before cursor is the partial path argument
		word := text
		if idx := strings.LastIndex(text, " "); idx >= 0 {
			word = text[idx+1:]
		}
		if isAbsolutePath(word) {
			return c.absolutePathSuggestions(word)
		}
		return c.relativePathSuggestions(word)
	}

	return nil
}

// PromptPrefix returns the currently active prompt prefix string.
func (c *Completer) PromptPrefix() string {
	return c.client.Name + " " + c.client.Pwd + "> "
}
