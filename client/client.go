package client

import (
	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/api"
	"strconv"
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
	Addr  string
	Token string
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

	return &Client{Vault: vault, Name: conf.Addr, Pwd: "", KVBackends: backends}, nil
}

// Read returns secret at given path, using given Client
func (client *Client) Read(path string) (secret *api.Secret, err error) {
	secret, err = client.Vault.Logical().Read(client.getKVDataPath(path))
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// Write writes secret to given path, using given Client
func (client *Client) Write(path string, secret *api.Secret) (err error) {
	_, err = client.Vault.Logical().Write(client.getKVDataPath(path), secret.Data)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes secret at given path, using given Client
func (client *Client) Delete(path string) (err error) {
	_, err = client.Vault.Logical().Delete(client.getKVMetaDataPath(path))
	if err != nil {
		return err
	}
	return nil
}

// List nodes at given path, using given Client
func (client *Client) List(path string) (result []string, err error) {
	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		return result, err
	}

	if s != nil {
		if keysInterface, ok := s.Data["keys"]; ok {
			for _, valInterface := range keysInterface.([]interface{}) {
				val := valInterface.(string)
				result = append(result, val)
			}
		}
	}

	return result, err
}

// IsFile checks if given path is a file
func (client *Client) IsFile(path string) (bool, error) {
	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		return false, err
	}

	if s == nil {
		return true, nil
	}
	return false, nil
}

// Traverse traverses paths via DFS and returns found paths array
func (client *Client) Traverse(path string) (paths []string) {
	var result []string
	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		log.Error("Error traversing path: %v", err)
		return nil
	}

	if s != nil {
		if keysInterface, ok := s.Data["keys"]; ok {
			for _, valInterface := range keysInterface.([]interface{}) {
				val := valInterface.(string)
				result = append(result, client.Traverse(path+"/"+val)...)
			}
		}
	} else {
		result = append(result, path)
	}

	return result
}
