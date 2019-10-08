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

var (
	vshVersion = ""
	verbose    = false
)

func printVersion() {
	fmt.Println(vshVersion)
}

func executor(in string) {
	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args := strings.Split(in, " ")
	commands := newCommands(vaultClient)
	var cmd cli.Command
	var run bool

	// Check for built-in commands.
	switch args[0] {
	case commands.ls.GetName():
		run = commands.ls.Parse(args)
		cmd = commands.ls
	case commands.cd.GetName():
		run = commands.cd.Parse(args)
		cmd = commands.cd
	case commands.mv.GetName():
		run = commands.mv.Parse(args)
		cmd = commands.mv
	case commands.cp.GetName():
		run = commands.cp.Parse(args)
		cmd = commands.cp
	case commands.rm.GetName():
		run = commands.rm.Parse(args)
		cmd = commands.rm
	case commands.cat.GetName():
		run = commands.cat.Parse(args)
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

	if run {
		if cmd.IsSane() {
			cmd.Run()
		}
	}
}

func main() {
	log.Init()

	var cmdString string
	var hasVersion bool
	flag.StringVar(&cmdString, "c", "", "command to run")
	flag.BoolVar(&hasVersion, "version", false, "print vsh version")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if hasVersion {
		printVersion()
		return
	}

	if verbose {
		log.ToggleVerbose()
	}

	conf := &client.VaultConfig{
		Addr:      os.Getenv("VAULT_ADDR"),
		Token:     os.Getenv("VAULT_TOKEN"),
		StartPath: os.Getenv("VAULT_PATH"),
	}
	var err error
	vaultClient, err = client.NewClient(conf)
	if err != nil {
		log.Error("Error initializing vault client | Are VAULT_ADDR, VAULT_TOKEN and VAULT_PATH properly set?")
		log.Error("%v", err)
		os.Exit(1)
	}

	if cmdString != "" {
		// Run non-interactive mode
		executor(cmdString)
	} else {
		// Run interactive mode
		completer := completer.NewCompleter(vaultClient)
		p := prompt.New(
			executor,
			completer.Complete,
			prompt.OptionTitle("vsh - interactive vault shell"),
			prompt.OptionLivePrefix(completer.PromptPrefix),
			prompt.OptionInputTextColor(prompt.Yellow),
			prompt.OptionShowCompletionAtStart(),
		)
		p.Run()
	}
}
