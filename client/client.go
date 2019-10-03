package client

import (
	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/api"
	"strconv"
	"strings"
)

// Client wrapper for Vault API client
type Client struct {
	Vault      *api.Client
	Name       string
	Pwd        string
	KVBackends map[string]int
}

// VaultConfig container to keep parameters for Client configuration
type VaultConfig struct {
	Addr      string
	Token     string
	StartPath string
}

// NewClient creates a new Client Vault wrapper
func NewClient(conf *VaultConfig) (*Client, error) {
	vault, err := api.NewClient(&api.Config{
		Address: conf.Addr,
	})

	if err != nil {
		return nil, err
	}

	vault.SetToken(conf.Token)

	mnts, err := vault.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var backends = make(map[string]int)
	for path, mnt := range mnts {
		if version, ok := mnt.Options["version"]; ok {
			v, err := strconv.Atoi(version)
			if err != nil {
				return nil, err
			}
			backends[path] = v
			log.Debug("Found KV backend '%v' with version '%v'", path, v)
		}
	}

	startPath := conf.StartPath
	if startPath == "" {
		startPath = "/"
	}

	if !strings.HasSuffix(startPath, "/") {
		startPath = startPath + "/"
	}

	return &Client{Vault: vault, Name: conf.Addr, Pwd: startPath, KVBackends: backends}, nil
}

// Read returns secret at given path, using given Client
func (client *Client) Read(absolutePath string) (secret *api.Secret, err error) {
	if client.isTopLevelPath(absolutePath) {
		secret, err = client.topLevelRead(absolutePath[1:])
	} else {
		secret, err = client.lowLevelRead(absolutePath[1:])
	}

	return secret, err
}

// Write writes secret to given path, using given Client
func (client *Client) Write(absolutePath string, secret *api.Secret) (err error) {
	if client.isTopLevelPath(absolutePath) {
		err = client.topLevelWrite(absolutePath[1:])
	} else {
		err = client.lowLevelWrite(absolutePath[1:], secret)
	}

	return err
}

// Delete deletes secret at given absolutePath, using given client
func (client *Client) Delete(absolutePath string) (err error) {
	if client.isTopLevelPath(absolutePath) {
		err = client.topLevelDelete(absolutePath[1:])
	} else {
		err = client.lowLevelDelete(absolutePath[1:])
	}

	return err
}

// List elements at the given absolutePath, using the given client
func (client *Client) List(absolutePath string) (result []string, err error) {
	if client.isTopLevelPath(absolutePath) {
		result = client.listTopLevel()
	} else {
		result, err = client.listLowLevel(absolutePath[1:])
	}

	return result, err
}

// GetType returns the file type the given absolutePath points to. Possible return values are BACKEND, NODE or LEAF
func (client *Client) GetType(absolutePath string) (kind PathKind, err error) {
	if client.isTopLevelPath(absolutePath) {
		kind, err = client.topLevelType()
	} else {
		kind, err = client.lowLevelType(absolutePath[1:])
	}

	return kind, err
}

// Traverse traverses given absolutePath via DFS and returns sub-paths in array
func (client *Client) Traverse(absolutePath string) (paths []string) {
	if client.isTopLevelPath(absolutePath) {
		paths = client.topLevelTraverse(absolutePath[1:])
	} else {
		paths = client.lowLevelTraverse(absolutePath[1:])
	}

	return paths
}
