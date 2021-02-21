package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/c-bata/go-prompt"
	"github.com/cosiner/argv"
	"github.com/fishi0x01/vsh/cli"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/completer"
	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/command/config"
)

var vaultClient *client.Client
var completerInstance *completer.Completer

var (
	vshVersion    = ""
	verbosity     = "INFO"
	isInteractive = true
)

type args struct {
	CmdString             string `arg:"-c,--cmd" help:"subcommand to run"`
	DisableAutoCompletion bool   `arg:"--disable-auto-completion" help:"disable auto-completion on paths"`
	Verbosity             string `arg:"-v,--log-level" help:"DEBUG | INFO | WARN | ERROR - debug option creates vsh_trace.log" default:"INFO" placeholder:"LEVEL"`
}

func (args) Version() string {
	return vshVersion
}

func (args) Description() string {
	return "vsh - Shell for Hashicorp Vault"
}

func executor(in string) {
	// Every command can change the vault content
	// i.e., the cache should be cleared on command execution
	vaultClient.ClearCache()

	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args, err := argv.Argv(in, func(backquoted string) (string, error) {
		return backquoted, nil
	}, nil)

	// edge cases
	if len(args) == 0 {
		fmt.Fprint(os.Stdout, "")
		if !isInteractive {
			os.Exit(1)
		}
		return
	}

	if err != nil {
		log.UserError("%v", err)
		return
	}
	commands := cli.NewCommands(vaultClient)
	var cmd cli.Command

	// parse command
	switch args[0][0] {
	case "toggle-auto-completion":
		completerInstance.TogglePathCompletion()
		return
	case "exit":
		os.Exit(0)
	default:
		cmd, err = getCommand(args[0], commands)
		if err == nil {
			err = cmd.Parse(args[0])
		}
	}

	if err != nil {
		log.UserError("%v", err)
		if cmd != nil {
			cmd.PrintUsage()
		}
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

func getCommand(args []string, commands *cli.Commands) (cmd cli.Command, err error) {
	cmd = commands.Get(args[0])
	if cmd == nil {
		return nil, fmt.Errorf("not a valid command: %s", args[0])
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
	var args args
	p := arg.MustParse(&args)
	switch level := strings.ToUpper(args.Verbosity); level {
	case "DEBUG", "INFO", "WARN", "ERROR":
		args.Verbosity = strings.ToUpper(args.Verbosity)
	default:
		p.Fail("Not a valid verbosity level")
	}

	var err error
	err = log.Init(args.Verbosity)
	if err != nil {
		os.Exit(1)
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

	if args.CmdString != "" {
		// Run non-interactive mode
		isInteractive = false
		executor(args.CmdString)
	} else {
		// Run interactive mode
		completerInstance = completer.NewCompleter(vaultClient, args.DisableAutoCompletion)
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
