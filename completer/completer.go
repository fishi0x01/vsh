package completer

import (
	"github.com/c-bata/go-prompt"
	"github.com/fishi0x01/vsh/client"
	"strings"
)

type Completer struct {
	client *client.Client
}

func NewCompleter(client *client.Client) (*Completer) {
	return &Completer{
		client: client,
	}
}

func isSupportedCommand(cmd string) bool {
	return cmd+" " == "cd " ||
		cmd+" " == "rm " ||
		cmd+" " == "mv " ||
		cmd+" " == "cat " ||
		cmd+" " == "ls "
}

func (c *Completer) getTopLevelSuggestions() []prompt.Suggest {
	var suggestions = []prompt.Suggest{}
	for k, _ := range c.client.KVBackends {
		suggestions = append(suggestions, prompt.Suggest{k, ""})
	}
	suggestions = append(suggestions, prompt.Suggest{".", ""})
	return suggestions
}

// Complete suggestions for completion
func (c *Completer) Complete(in prompt.Document) []prompt.Suggest {
	com := strings.Split(in.TextBeforeCursor(), " ")[0]
	cur := in.GetWordBeforeCursor()

	var suggestions = []prompt.Suggest{}
	if (isSupportedCommand(com)) {
		if (c.client.Pwd == "") {
			suggestions = c.getTopLevelSuggestions()
		} else {
			completePath := c.client.Pwd + cur
			li := strings.LastIndex(completePath, "/")
			if (li > 0) {
				completePath = completePath[:li+1]
			} else {
				completePath = completePath[0:li]
			}
			options, err := c.client.List(completePath)
			if err != nil {
				// TODO: handle error
			}
			options = append(options, ".", "..")
			for _, node := range options {
				suggestions = append(suggestions, prompt.Suggest{node, ""})
			}
		}
	}
	return prompt.FilterHasPrefix(suggestions, cur, true)
}

// PromptPrefix returns the currently active prompt prefix
func (c *Completer) PromptPrefix() (string, bool) {
	var p string
	parts := strings.Split(c.client.Pwd, "/")
	if (len(parts) > 1) {
		p = parts[len(parts)-2] + "/"
	} else {
		p = parts[0]
	}
	return c.client.Name + " " + p + "> ", true
}
