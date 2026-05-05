package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cosiner/argv"
	"github.com/fishi0x01/vsh/internal/cli"
	"github.com/fishi0x01/vsh/internal/client"
	"github.com/fishi0x01/vsh/internal/completer"
	"github.com/fishi0x01/vsh/internal/logger"
	"github.com/fishi0x01/vsh/internal/tui"
	"github.com/hashicorp/vault/api/cliconfig"
)

var vshVersion = ""

type app struct {
	vaultClient   *client.Client
	completer     *completer.Completer
	workerCount   int
	isInteractive bool
}

type args struct {
	CmdString             string `arg:"-c,--cmd"                  help:"subcommand to run"`
	DisableAutoCompletion bool   `arg:"--disable-auto-completion" help:"disable auto-completion on paths"`
	Verbosity             string `arg:"-v,--log-level"            help:"DEBUG | INFO | WARN | ERROR - debug option creates vsh_trace.log" default:"INFO" placeholder:"LEVEL"`
	WorkerCount           int    `arg:"--worker-count"            help:"concurrent workers for recursive ops"                             default:"10"`
}

func (args) Version() string {
	return vshVersion
}

func (args) Description() string {
	return "vsh - Shell for Hashicorp Vault"
}

func (a *app) executor(in string) {
	// Every command can change the vault content
	// i.e., the cache should be cleared after a command got executed
	defer a.vaultClient.ClearCache()

	// Split the input separate the command and the arguments.
	in = strings.TrimSpace(in)
	args, err := argv.Argv(in, func(backquoted string) (string, error) {
		return backquoted, nil
	}, nil)

	// edge cases
	if len(args) == 0 {
		_, err := fmt.Fprint(os.Stdout, "")
		if err != nil {
			logger.UserError("Error printing: %v", err)
		}
		if !a.isInteractive {
			os.Exit(1)
		}
		return
	}

	if err != nil {
		logger.UserError("%v", err)
		return
	}
	commands := cli.NewCommands(a.vaultClient, a.workerCount, a.isInteractive)
	var cmd cli.Command

	// parse command
	switch args[0][0] {
	case "toggle-auto-completion":
		a.completer.TogglePathCompletion()
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
		logger.UserError("%v", err)
		if cmd != nil {
			cmd.PrintUsage()
		}
	}

	if err != nil && !a.isInteractive {
		os.Exit(1)
	}

	if err == nil && cmd.IsSane() {
		if a.isInteractive {
			if sa, ok := cmd.(cli.SpinnerAware); ok {
				sa.SetSpinnerControl(tui.PauseSpinner, tui.ResumeSpinner)
			}
		}
		ret := cmd.Run()
		if !a.isInteractive {
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

// parseOnly parses the command line and returns the Command without running it.
// Returns nil if parsing fails or the input is a built-in (exit, toggle-...).
// Used to check for Confirmer before starting RunWithSpinner.
func (a *app) parseOnly(in string) cli.Command {
	in = strings.TrimSpace(in)
	cmdArgs, err := argv.Argv(in, func(s string) (string, error) { return s, nil }, nil)
	if err != nil || len(cmdArgs) == 0 {
		return nil
	}
	commands := cli.NewCommands(a.vaultClient, a.workerCount, a.isInteractive)
	cmd, err := getCommand(cmdArgs[0], commands)
	if err != nil {
		return nil
	}
	if err = cmd.Parse(cmdArgs[0]); err != nil {
		return nil
	}
	return cmd
}

func getVaultToken() (token string, err error) {
	token = os.Getenv("VAULT_TOKEN")
	if token == "" {
		helper, err := cliconfig.DefaultTokenHelper()
		if err != nil {
			return "", err
		}
		tok, err := helper.Get()
		if err != nil {
			return "", err
		}
		token = tok
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
	err = logger.Init(args.Verbosity)
	if err != nil {
		os.Exit(1)
	}
	defer logger.Close()

	token, ve := getVaultToken()
	if ve != nil {
		logger.AppError("Error getting vault token")
		logger.AppError("%v", ve)
		return
	}

	conf := &client.VaultConfig{
		Addr:            os.Getenv("VAULT_ADDR"),
		Token:           token,
		StartPath:       os.Getenv("VAULT_PATH"),
		CertificatePath: os.Getenv("VAULT_CACERT"),
	}

	a := &app{
		workerCount:   args.WorkerCount,
		isInteractive: true,
	}

	a.vaultClient, err = client.NewClient(conf)
	if err != nil {
		logger.UserError(
			"Error initializing vault client | Is VAULT_ADDR properly set? Do you provide a proper token?",
		)
		logger.UserError("%v", err)
		os.Exit(1)
	}

	if args.CmdString != "" {
		// Run non-interactive mode
		a.isInteractive = false
		a.executor(args.CmdString)
	} else {
		// Run interactive mode — quit/restart loop so executor output reaches
		// the terminal freely between bubbletea sessions.
		a.completer = completer.NewCompleter(a.vaultClient, args.DisableAutoCompletion)
		var history []string
		for {
			m := tui.NewModel(a.completer.Complete, a.completer.PromptPrefix, history)
			result, err := tea.NewProgram(m).Run()
			if err != nil {
				logger.AppError("Error running TUI: %v", err)
				os.Exit(1)
			}
			final := result.(tui.Model)
			history = final.History()
			if final.Quitting() {
				break
			}
			cmd := final.LastCommand()
			fmt.Printf("%s%s\n", tui.RenderPrefix(a.completer.PromptPrefix()), cmd)
			if cmd != "" {
				parsed := a.parseOnly(cmd)
				if c, ok := parsed.(cli.Confirmer); ok && !c.Confirm() {
					a.vaultClient.ClearCache()
				} else {
					tui.RunWithSpinner(func() { a.executor(cmd) })
				}
			}
		}
	}
}
