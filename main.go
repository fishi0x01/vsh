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

var (
	vshVersion    = ""
	verbosity     = "INFO"
	isInteractive = true
)

func printVersion() {
	fmt.Println(vshVersion)
}

func parseInput(line string) (args []string) {
	quote := '0'
	escaped := false
	arg := ""

	for _, c := range line {
		switch {
		case c == '\\' && !escaped: // next char will be escaped
			escaped = true
			continue
		case escaped:
			escaped = false
			arg += string(c) // append char to current arg buffer
			continue
		case c == quote: // terminating quote
			quote = '0'
			args = append(args, arg)
			arg = ""
		case c == '"' || c == '\'':
			if quote == '0' { // beginning quote
				quote = c
			} else if c != quote { // non-matching quote char
				arg += string(c)
			}
		case c == ' ':
			if quote == '0' {
				if arg != "" { // unquoted space, store non-empty arg
					args = append(args, arg)
					arg = ""
				}
				continue
			}
			fallthrough
		default:
			arg += string(c) // append char to current arg buffer
		}
	}

	if arg != "" { // store non-empty arg
		args = append(args, arg)
	}
	return args
}

var completerInstance *completer.Completer

func executor(in string) {
	// Every command can change the vault content
	// i.e., the cache should be cleared on command execution
	vaultClient.ClearCache()

	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args := parseInput(in)
	commands := cli.NewCommands(vaultClient)
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
		log.UserError("%v", err)
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

func getCommand(args []string, commands *cli.Commands) (cmd cli.Command, err error) {
	switch args[0] {
	case commands.Ls.GetName():
		err = commands.Ls.Parse(args)
		cmd = commands.Ls
	case commands.Cd.GetName():
		err = commands.Cd.Parse(args)
		cmd = commands.Cd
	case commands.Mv.GetName():
		err = commands.Mv.Parse(args)
		cmd = commands.Mv
	case commands.Append.GetName():
		err = commands.Append.Parse(args)
		cmd = commands.Append
	case commands.Cp.GetName():
		err = commands.Cp.Parse(args)
		cmd = commands.Cp
	case commands.Rm.GetName():
		err = commands.Rm.Parse(args)
		cmd = commands.Rm
	case commands.Cat.GetName():
		err = commands.Cat.Parse(args)
		cmd = commands.Cat
	case commands.Grep.GetName():
		err = commands.Grep.Parse(args)
		cmd = commands.Grep
	case commands.Replace.GetName():
		err = commands.Replace.Parse(args)
		cmd = commands.Replace
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
