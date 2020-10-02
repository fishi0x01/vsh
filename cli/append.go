package cli

import (
	"fmt"
	"io"
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

	client *client.Client
	stderr io.Writer
	stdout io.Writer
	Source string
	Target string
	Mode   AppendMode
}

// NewAppendCommand creates a new AppendCommand parameter container
func NewAppendCommand(c *client.Client, stdout io.Writer, stderr io.Writer) *AppendCommand {
	return &AppendCommand{
		name:   "append",
		client: c,
		stdout: stdout,
		stderr: stderr,
		Mode:   ModeSkip,
	}
}

// GetName returns the AppendCommand's name identifier
func (cmd *AppendCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *AppendCommand) IsSane() bool {
	return cmd.Source != "" && cmd.Target != "" && cmd.Mode != ModeInvalid
}

func printUsage() {
	fmt.Println("Usage:\nappend <from> <to> [-f|--force|-r|--rename|-s|--skip]")
}

func isFlag(flag string) bool {
	return strings.HasPrefix(flag, "-")
}

func parseFlag(flag string) AppendMode {
	switch strings.TrimSpace(flag) {
	case "-f", "--force":
		return ModeOverwrite
	case "", "-s", "--skip":
		return ModeSkip
	case "-r", "--rename":
		return ModeRename
	default:
		return ModeInvalid
	}
}

func (cmd *AppendCommand) parseArgs(src, dest, flag string) bool {
	cmd.Source = src
	cmd.Target = dest
	mode := parseFlag(flag)
	cmd.Mode = mode
	if mode == ModeInvalid {
		return false
	}
	return true
}

// tryParse returns true when parsing succeeded, false otherwise
func (cmd *AppendCommand) tryParse(args []string) (success bool) {
	if len(args) == 3 {
		return cmd.parseArgs(args[1], args[2], "--skip") // --skip is default
	}
	if len(args) == 4 {
		// flag can be given at the end or immediately after `append`
		if isFlag(args[3]) {
			return cmd.parseArgs(args[1], args[2], args[3])
		}
		if isFlag(args[1]) {
			return cmd.parseArgs(args[2], args[3], args[1])
		}
	}
	// wrong number of params or flag at incorrect position
	return false
}

// Parse parses the arguments and returns true on success; otherwise it prints usage and returns false
func (cmd *AppendCommand) Parse(args []string) error {
	success := cmd.tryParse(args)
	if !success {
		printUsage()
		return fmt.Errorf("cannot parse arguments")
	}
	return nil
}

// Run executes 'append' with given AppendCommand's parameters
func (cmd *AppendCommand) Run() int {
	newSrcPwd := cmdPath(cmd.client.Pwd, cmd.Source)
	newTargetPwd := cmdPath(cmd.client.Pwd, cmd.Target)

	src := cmd.client.GetType(newSrcPwd)
	if src != client.LEAF {
		log.NotAValidPath(newSrcPwd)
		return 1
	}

	if err := cmd.mergeSecrets(newSrcPwd, newTargetPwd); err != nil {
		log.Error("Append failed: " + err.Error())
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
	if targetSecret == nil {
		if err = cmd.client.Write(target, &api.Secret{Data: dummy}); err != nil {
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
	merged := make(map[string]interface{})
	skippedKeys := make([]string, 0)

	for k, v := range targetSecret.Data {
		if rec, ok := v.(map[string]interface{}); ok {
			for kk, vv := range rec {
				merged[kk] = vv
			}
		} else {
			merged[k] = v
		}
	}

	for k, v := range sourceSecret.Data {
		if rec, ok := v.(map[string]interface{}); ok {
			for kk, vv := range rec {
				skipped := addKey(merged, onConflict, kk, vv)
				skippedKeys = append(skippedKeys, skipped...)
			}
		} else {
			skipped := addKey(merged, onConflict, k, v)
			skippedKeys = append(skippedKeys, skipped...)
		}
	}
	// write
	if err := cmd.client.Write(target, &api.Secret{Data: merged}); err != nil {
		fmt.Println(err)
		return err
	}
	log.Info("Appended values from %s to %s", source, target)
	if len(skippedKeys) > 0 {
		log.Info("Handled conflicting keys according to the '%s' strategy. Keys: %s", onConflict, strings.Join(skippedKeys, ", "))
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
