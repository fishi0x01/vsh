package cli

import (
	"github.com/fishi0x01/vsh/client"
	"path/filepath"
	"strings"
)

// Command interface to describe a command structure
type Command interface {
	Run() error
	GetName() string
}

func cmdPath(pwd string, arg string) (result string) {
	result =filepath.Clean(pwd + arg)
	if strings.HasPrefix(arg, "/") {
		// absolute path is given
		result = arg
	}
	return result
}

func runCommandWithTraverseTwoPaths(client *client.Client, source string, target string, f func(string, string) error) error {
	for _, path := range client.Traverse(source) {
		target := strings.Replace(path, source, target, 1)
		err := f(path, target)
		if err != nil {
			return err
		}
	}

	return nil
}
