package cli

import "strings"

// Command interface to describe a command structure
type Command interface {
	Run() error
	GetName() string
}

func cmdPath(pwd string, arg string) (result string) {
	result = pwd + arg
	if strings.HasPrefix(arg, "/") {
		// absolute path is given
		result = arg
	}
	return result
}
