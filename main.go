package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fishi0x01/vsh/cli"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/completer"
	"github.com/fishi0x01/vsh/log"
	"os"
	"strings"
)

var vaultClient *client.Client

type commands struct {
	mv  *cli.MoveCommand
	cp  *cli.CopyCommand
	rm  *cli.RemoveCommand
	ls  *cli.ListCommand
	cd  *cli.CdCommand
	cat *cli.CatCommand
}

func newCommands(client *client.Client) *commands {
	return &commands{
		mv:  cli.NewMoveCommand(client, os.Stdout, os.Stderr),
		cp:  cli.NewCopyCommand(client, os.Stdout, os.Stderr),
		rm:  cli.NewRemoveCommand(client, os.Stdout, os.Stderr),
		ls:  cli.NewListCommand(client, os.Stdout, os.Stderr),
		cd:  cli.NewCdCommand(client, os.Stdout, os.Stderr),
		cat: cli.NewCatCommand(client, os.Stdout, os.Stderr),
	}
}

func executor(in string) {
	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args := strings.Split(in, " ")
	commands := newCommands(vaultClient)
	var cmd cli.Command

	// Check for built-in commands.
	switch args[0] {
	case commands.ls.GetName():
		// 'ls' the current path
		if len(args) > 1 {
			commands.ls.Path = args[1]
		} else {
			commands.ls.Path = vaultClient.Pwd
		}
		cmd = commands.ls
	case commands.cd.GetName():
		// 'cd' to path
		commands.cd.Path = args[1]
		cmd = commands.cd
	case commands.mv.GetName():
		// 'mv' the current path
		commands.mv.Source = args[1]
		commands.mv.Target = args[2]
		cmd = commands.mv
	case commands.cp.GetName():
		// 'cp' the current path
		commands.cp.Source = args[1]
		commands.cp.Target = args[2]
		cmd = commands.cp
	case commands.rm.GetName():
		// 'rm' the current path
		commands.rm.Path = args[1]
		cmd = commands.rm
	case commands.cat.GetName():
		// 'cat' given file
		commands.cat.Path = args[1]
		cmd = commands.cat
	case "exit":
		os.Exit(0)
	case "":
		fmt.Fprint(os.Stdout, "")
		return
	default:
		fmt.Fprintln(os.Stderr, "Unknown command '"+args[0]+"'")
		return
	}

	cmd.Run()
}

func main() {
	log.Init()

	var inputString string
	flag.StringVar(&inputString, "c", "", "command to run")
	flag.Parse()

	conf := &client.VaultConfig{
		Addr: os.Getenv("VAULT_ADDR"),
		Token: os.Getenv("VAULT_TOKEN"),
		StartPath: os.Getenv("VAULT_PATH"),
	}
	var err error
	vaultClient, err = client.NewClient(conf)
	if err != nil {
		log.Error("Error initializing vault client | Are VAULT_ADDR and VAULT_TOKEN properly set?")
		log.Error("%v", err)
		os.Exit(1)
	}

	if inputString != "" {
		// Run non-interactive mode
		executor(inputString)
	} else {
		// Run interactive mode
		completer := completer.NewCompleter(vaultClient)
		p := prompt.New(
			executor,
			completer.Complete,
			prompt.OptionTitle("vsh - interactive vault shell"),
			prompt.OptionLivePrefix(completer.PromptPrefix),
			prompt.OptionInputTextColor(prompt.Yellow),
		)
		p.Run()
	}
}
