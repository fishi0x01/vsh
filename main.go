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
		mv:     cli.NewMoveCommand(client),
		cp:     cli.NewCopyCommand(client),
		append: cli.NewAppendCommand(client),
		rm:     cli.NewRemoveCommand(client),
		ls:     cli.NewListCommand(client),
		cd:     cli.NewCdCommand(client),
		cat:    cli.NewCatCommand(client),
		grep:   cli.NewGrepCommand(client, os.Stdout, os.Stderr),
	}
}

var (
	vshVersion    = ""
	verbosity     = "INFO"
	isInteractive = true
)

func printVersion() {
	fmt.Println(vshVersion)
}

func parseInput(line string) (args []string) {
	// TODO: allow "" and "\"\""
	return strings.Fields(line)
}

var completerInstance *completer.Completer

func executor(in string) {
	// Every command can change the vault content
	// i.e., the cache should be cleared on command execution
	vaultClient.ClearCache()

	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args := parseInput(in)
	commands := newCommands(vaultClient)
	var cmd cli.Command
	var err error

	// edge cases
	if len(args) == 0 {
		fmt.Fprint(os.Stdout, "")
		if !isInteractive {
			os.Exit(1)
		}
		return
	}

	// parse command
	switch args[0] {
	case "toggle-auto-completion":
		completerInstance.TogglePathCompletion()
		return
	case "exit":
		os.Exit(0)
	default:
		cmd, err = getCommand(args, commands)
	}

	if err != nil && cmd != nil {
		cmd.PrintUsage()
	}

	if err != nil && !isInteractive {
		os.Exit(1)
	}

	if err == nil && cmd.IsSane() {
		ret := cmd.Run()
		if !isInteractive {
			os.Exit(ret)
		}
	}
}

func getCommand(args []string, commands *commands) (cmd cli.Command, err error) {
	switch args[0] {
	case commands.ls.GetName():
		err = commands.ls.Parse(args)
		cmd = commands.ls
	case commands.cd.GetName():
		err = commands.cd.Parse(args)
		cmd = commands.cd
	case commands.mv.GetName():
		err = commands.mv.Parse(args)
		cmd = commands.mv
	case commands.append.GetName():
		err = commands.append.Parse(args)
		cmd = commands.append
	case commands.cp.GetName():
		err = commands.cp.Parse(args)
		cmd = commands.cp
	case commands.rm.GetName():
		err = commands.rm.Parse(args)
		cmd = commands.rm
	case commands.cat.GetName():
		err = commands.cat.Parse(args)
		cmd = commands.cat
	case commands.grep.GetName():
		err = commands.grep.Parse(args)
		cmd = commands.grep
	default:
		log.UserError("Not a valid command: %s", args[0])
		return nil, fmt.Errorf("not a valid command")
	}
	return cmd, err
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
	var cmdString string
	var hasVersion bool
	var disableAutoCompletion bool
	flag.StringVar(&cmdString, "c", "", "command to run")
	flag.BoolVar(&hasVersion, "version", false, "print vsh version")
	flag.BoolVar(&disableAutoCompletion, "disable-auto-completion", false, "disable auto-completion on paths")
	flag.StringVar(&verbosity, "v", "INFO", "DEBUG | INFO | WARN | ERROR - debug option creates vsh_trace.log")
	flag.Parse()

	var err error
	err = log.Init(verbosity)
	if err != nil {
		os.Exit(1)
	}

	if hasVersion {
		printVersion()
		return
	}

	token, ve := getVaultToken()
	if ve != nil {
		log.AppError("Error getting vault token")
		log.AppError("%v", ve)
		return
	}

	conf := &client.VaultConfig{
		Addr:      os.Getenv("VAULT_ADDR"),
		Token:     token,
		StartPath: os.Getenv("VAULT_PATH"),
	}

	vaultClient, err = client.NewClient(conf)
	if err != nil {
		log.UserError("Error initializing vault client | Is VAULT_ADDR properly set? Do you provide a proper token?")
		log.UserError("%v", err)
		os.Exit(1)
	}

	if cmdString != "" {
		// Run non-interactive mode
		isInteractive = false
		executor(cmdString)
	} else {
		// Run interactive mode
		completerInstance = completer.NewCompleter(vaultClient, disableAutoCompletion)
		p := prompt.New(
			executor,
			completerInstance.Complete,
			prompt.OptionTitle("vsh - interactive vault shell"),
			prompt.OptionLivePrefix(completerInstance.PromptPrefix),
			prompt.OptionInputTextColor(prompt.Yellow),
			prompt.OptionShowCompletionAtStart(),
		)
		p.Run()
	}
}
