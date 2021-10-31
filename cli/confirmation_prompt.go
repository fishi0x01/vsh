package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Inspiration from: https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func askForConfirmation(s string) (bool, error) {
	c := color.New(color.FgCyan)
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < 2; i++ {
		c.Printf("%s [y/N]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("Error reading input")
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}

	return false, nil
}
