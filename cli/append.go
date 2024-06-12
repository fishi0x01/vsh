package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/api"
)

// AppendMode defines behaviour in case existing secrets are in conflict
type AppendMode string

const (
	// ModeSkip do not append secret on conflict
	ModeSkip AppendMode = "skip"
	// ModeOverwrite overwrite existing secret on conflict
	ModeOverwrite AppendMode = "overwrite"
	// ModeRename keep existing secret and create a new key on conflict
	ModeRename AppendMode = "rename"
	// ModeInvalid denotes invalid mode
	ModeInvalid AppendMode = "invalid"
)

// AppendCommand container for all 'append' parameters
type AppendCommand struct {
	name string
	args *AppendCommandArgs

	client *client.Client
	Mode   AppendMode
}

// AppendCommandArgs provides a struct for go-arg parsing
type AppendCommandArgs struct {
	Source string `arg:"positional,required"`
	Target string `arg:"positional,required"`
	Force  bool   `arg:"-f,--force" help:"Overwrite key if exists"`
	Skip   bool   `arg:"-s,--skip" help:"Skip key if exists (default)"`
	Rename bool   `arg:"-r,--rename" help:"Rename key if exists"`
}

// Description provides detail on what the command does
func (AppendCommandArgs) Description() string {
	return "appends the contents of one secret to another"
}

// NewAppendCommand creates a new AppendCommand parameter container
func NewAppendCommand(c *client.Client) *AppendCommand {
	return &AppendCommand{
		name:   "append",
		client: c,
		args:   &AppendCommandArgs{},
		Mode:   ModeSkip,
	}
}

// GetName returns the AppendCommand's name identifier
func (cmd *AppendCommand) GetName() string {
	return cmd.name
}

// GetArgs provides the struct holding arguments for the command
func (cmd *AppendCommand) GetArgs() interface{} {
	return cmd.args
}

// IsSane returns true if command is sane
func (cmd *AppendCommand) IsSane() bool {
	return cmd.args.Source != "" && cmd.args.Target != "" && cmd.Mode != ModeInvalid
}

// PrintUsage print command usage
func (cmd *AppendCommand) PrintUsage() {
	fmt.Println(Help(cmd))
}

// Parse parses the arguments into the Command and Args structs
func (cmd *AppendCommand) Parse(args []string) error {
	_, err := parseCommandArgs(args, cmd)
	if err != nil {
		return err
	}

	if cmd.args.Skip == true {
		cmd.Mode = ModeSkip
	} else if cmd.args.Force == true {
		cmd.Mode = ModeOverwrite
	} else if cmd.args.Rename == true {
		cmd.Mode = ModeRename
	}
	return nil
}

// Run executes 'append' with given AppendCommand's parameters
func (cmd *AppendCommand) Run() int {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.args.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.args.Target)

	src := cmd.client.GetType(newSrcPwd)
	if src != client.LEAF {
		log.UserError("Not a valid path for operation: %s", newSrcPwd)
		return 1
	}

	if err := cmd.mergeSecrets(newSrcPwd, newTargetPwd); err != nil {
		log.AppError("Append failed: " + err.Error())
		return 1
	}
	return 0
}

func (cmd *AppendCommand) createDummySecret(target string) error {
	targetSecret, err := cmd.client.Read(target)
	if targetSecret != nil && err == nil {
		return nil
	}

	dummy := make(map[string]interface{})
	dummy["placeholder"] = struct{}{}
	dummySecret := client.NewSecret(&api.Secret{Data: dummy}, target)
	if targetSecret == nil {
		if err = cmd.client.Write(target, dummySecret); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *AppendCommand) mergeSecrets(source string, target string) error {
	sourceSecret, err := cmd.client.Read(source)
	if err != nil {
		return err
	}
	cmd.createDummySecret(target)
	targetSecret, err := cmd.client.Read(target)
	if err != nil {
		return err
	}

	onConflict := cmd.Mode
	merged := targetSecret.GetData()
	skippedKeys := make([]string, 0)

	for k, v := range sourceSecret.GetData() {
		skipped := addKey(merged, onConflict, k, v)
		skippedKeys = append(skippedKeys, skipped...)
	}

	// write
	resultSecret := client.NewSecret(&api.Secret{Data: merged}, target)
	if err := cmd.client.Write(target, resultSecret); err != nil {
		fmt.Println(err)
		return err
	}
	log.UserDebug("Appended values from %s to %s", source, target)
	if len(skippedKeys) > 0 {
		log.UserDebug("Handled conflicting keys according to the '%s' strategy. Keys: %s", onConflict, strings.Join(skippedKeys, ", "))
	}
	return nil
}

func addKey(merged map[string]interface{}, onConflict AppendMode, key string, value interface{}) []string {
	skippedKeys := make([]string, 0)
	// if this key is already present in the destination
	if _, ok := merged[key]; ok {
		switch onConflict {
		case ModeOverwrite:
			merged[key] = value
		case ModeSkip, ModeInvalid:
			skippedKeys = append(skippedKeys, key)
		case ModeRename:
			key2 := getNextFreeKey(merged, key)
			merged[key2] = value
		}
	} else {
		merged[key] = value
	}
	return skippedKeys
}

func getNextFreeKey(merged map[string]interface{}, key string) string {
	idx := 0
	ok := true // assume key exists
	candidate := key
	for ok {
		idx++
		candidate = key + "_" + strconv.Itoa(idx)
		_, ok = merged[candidate]
	}
	return candidate
}
