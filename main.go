package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fishi0x01/vsh/cli"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/completer"
	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/command/config"
)

var vaultClient *client.Client

type commands struct {
	mv     *cli.MoveCommand
	cp     *cli.CopyCommand
	append *cli.AppendCommand
	rm     *cli.RemoveCommand
	ls     *cli.ListCommand
	cd     *cli.CdCommand
	cat    *cli.CatCommand
	grep   *cli.GrepCommand
}

func newCommands(client *client.Client) *commands {
	return &commands{
		mv:     cli.NewMoveCommand(client, os.Stdout, os.Stderr),
		cp:     cli.NewCopyCommand(client, os.Stdout, os.Stderr),
		append: cli.NewAppendCommand(client, os.Stdout, os.Stderr),
		rm:     cli.NewRemoveCommand(client, os.Stdout, os.Stderr),
		ls:     cli.NewListCommand(client, os.Stdout, os.Stderr),
		cd:     cli.NewCdCommand(client, os.Stdout, os.Stderr),
		cat:    cli.NewCatCommand(client, os.Stdout, os.Stderr),
		grep:   cli.NewGrepCommand(client, os.Stdout, os.Stderr),
	}
}

var (
	vshVersion = ""
	verbose    = false
)

func printVersion() {
	fmt.Println(vshVersion)
}

func parseInput(line string) (args []string) {
	// TODO: allow "" and "\"\""
	return strings.Fields(line)
}

func executor(in string) {
	// Every command can change the vault content
	// i.e., the cache should be cleared on command execution
	vaultClient.ClearCache()

	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args := parseInput(in)
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
	case commands.append.GetName():
		run = commands.append.Parse(args)
		cmd = commands.append
	case commands.cp.GetName():
		run = commands.cp.Parse(args)
		cmd = commands.cp
	case commands.rm.GetName():
		run = commands.rm.Parse(args)
		cmd = commands.rm
	case commands.cat.GetName():
		run = commands.cat.Parse(args)
		cmd = commands.cat
	case commands.grep.GetName():
		run = commands.grep.Parse(args)
		cmd = commands.grep
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

func getVaultToken() (token string, err error) {
	token = os.Getenv("VAULT_TOKEN")
	if token == "" {
		helper, ve := config.DefaultTokenHelper()
		if ve != nil {
			err = ve
			return token, err
		}
		token, ve = helper.Get()
		if ve != nil {
			err = ve
		}
	}
	return token, err
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

	token, ve := getVaultToken()
	if ve != nil {
		log.Error("Error getting vault token")
		log.Error("%v", ve)
		return
	}

	conf := &client.VaultConfig{
		Addr:      os.Getenv("VAULT_ADDR"),
		Token:     token,
		StartPath: os.Getenv("VAULT_PATH"),
	}
	var err error
	vaultClient, err = client.NewClient(conf)
	if err != nil {
		log.Error("Error initializing vault client | Is VAULT_ADDR properly set? Do you provide a proper token?")
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
