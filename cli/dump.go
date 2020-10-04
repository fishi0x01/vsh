package cli

import (
	"crypto/sha1"
	"fmt"
	"github.com/fishi0x01/vsh/client"
	"github.com/fishi0x01/vsh/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DumpCommand container for all 'dump' parameters
type DumpCommand struct {
	name string

	client *client.Client
	Path   string
}

// NewDumpCommand creates a new DumpCommand parameter container
func NewDumpCommand(c *client.Client) *DumpCommand {
	return &DumpCommand{
		name:   "dump",
		client: c,
	}
}

// GetName returns the DumpCommand's name identifier
func (cmd *DumpCommand) GetName() string {
	return cmd.name
}

// IsSane returns true if command is sane
func (cmd *DumpCommand) IsSane() bool {
	return cmd.Path != ""
}

// PrintUsage print command usage
func (cmd *DumpCommand) PrintUsage() {
	log.UserInfo("Usage:\ndump <path>")
}

// Parse given arguments and return status
func (cmd *DumpCommand) Parse(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("cannot parse arguments")
	}
	cmd.Path = args[1]
	return nil
}

// Run executes 'dump' with given DumpCommand's parameters
func (cmd *DumpCommand) Run() int {
	newPwd := cmdPath(cmd.client.Pwd, cmd.Path)
	rootDir := "vsh-dump_" + time.Now().Format("2006_0102_150405")

	switch t := cmd.client.GetType(newPwd); t {
	case client.LEAF:
		cmd.dumpSecret(newPwd, rootDir)
	case client.NODE:
		for _, path := range cmd.client.Traverse(newPwd) {
			err := cmd.dumpSecret(path, rootDir)
			if err != nil {
				return 1
			}
		}
	default:
		log.UserError("Not a valid path for operation: %s", newPwd)
		return 1
	}

	return 0
}

func (cmd *DumpCommand) dumpSecret(path string, rootDir string) error {
	log.UserDebug("Dump %s", path)
	secret, err := cmd.client.Read(path)
	if err != nil {
		return err
	}

	// use hashes to prevent ambiguous file/dir
	h := sha1.New()
	h.Write([]byte(filepath.Dir(path)))
	relativeDumpFilePath := filepath.Join(fmt.Sprintf("%x", h.Sum(nil)), filepath.Base(path))
	dumpFilePath := filepath.Join(rootDir, relativeDumpFilePath)
	os.MkdirAll(filepath.Dir(dumpFilePath), os.ModePerm)
	dumpFile, err := os.OpenFile(dumpFilePath, os.O_CREATE|os.O_WRONLY, 0400)
	if err != nil {
		log.AppError("%+v", err)
		return err
	}
	defer dumpFile.Close()
	jsonContent := "{"
	isFirstValue := true

	for k, v := range secret.Data {
		if !isFirstValue {
			jsonContent += ","
		}
		isFirstValue = false
		if rec, ok := v.(map[string]interface{}); ok {
			// KV 2
			isFirstValue = true
			for kk, vv := range rec {
				if !isFirstValue {
					jsonContent += ","
				}
				val := fmt.Sprintf("\"%s\":\"%s\"", kk, vv)
				val = strings.ReplaceAll(val, "\n", `\n`)
				jsonContent += val
				isFirstValue = false
			}
			isFirstValue = false
		} else {
			// KV 1
			jsonContent += fmt.Sprintf("\"%s\":\"%s\"", k, v)
		}
	}
	jsonContent += "}"
	dumpFile.WriteString(jsonContent)
	if err != nil {
		log.AppError("%+v", err)
		return err
	}

	scriptFilePath := filepath.Join(rootDir, "restore.sh")
	scriptFile, err := os.OpenFile(scriptFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.AppError("%+v", err)
		return err
	}
	defer scriptFile.Close()
	scriptFile.WriteString(fmt.Sprintf("vault kv put %s @%s\n", strings.TrimPrefix(path, "/"), relativeDumpFilePath))
	log.UserDebug("vault kv put %s @%s", strings.TrimPrefix(path, "/"), relativeDumpFilePath)

	return nil
}
