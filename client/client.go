package client

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/api"
)

// Client wrapper for Vault API client
type Client struct {
	Vault      *api.Client
	Name       string
	Pwd        string
	KVBackends map[string]int
	listCache  map[string][]string
	cacheMutex sync.Mutex
	waitGroup  sync.WaitGroup
}

// VaultConfig container to keep parameters for Client configuration
type VaultConfig struct {
	Addr      string
	Token     string
	StartPath string
}

type secretOperation struct {
	Result *Secret
	Path   string
	Error  error
}

func verifyClientPwd(client *Client) (*Client, error) {
	if client.Pwd == "" {
		client.Pwd = "/"
	}

	if !strings.HasSuffix(client.Pwd, "/") {
		client.Pwd = client.Pwd + "/"
	}

	if !strings.HasPrefix(client.Pwd, "/") {
		client.Pwd = "/" + client.Pwd
	}

	t := client.GetType(client.Pwd)

	if t != NODE && t != BACKEND {
		return nil, errors.New("VAULT_PATH is not a valid directory path")
	}

	return client, nil
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

	permissions, err := vault.Sys().CapabilitiesSelf("sys/mounts")

	var mounts map[string]*api.MountOutput
	if sliceContains(permissions, "list") || sliceContains(permissions, "root") {
		mounts, err = vault.Sys().ListMounts()
	} else {
		log.UserDebug("Cannot auto-discover mount backends: Token does not have list permission on sys/mounts")
	}

	if err != nil {
		log.AppTrace("%+v", err)
		return nil, err
	}

	var backends = make(map[string]int)
	for path, mount := range mounts {
		if version, ok := mount.Options["version"]; ok {
			v, err := strconv.Atoi(version)
			if err != nil {
				return nil, err
			}
			backends[path] = v
			log.UserDebug("Found KV backend '%v' with version '%v'", path, v)
		}
	}

	return verifyClientPwd(&Client{
		Vault:      vault,
		Name:       conf.Addr,
		Pwd:        conf.StartPath,
		KVBackends: backends,
		listCache:  make(map[string][]string),
	})
}

// Read returns secret at given path, using given Client
func (client *Client) Read(absolutePath string) (secret *Secret, err error) {
	var apiSecret *api.Secret
	if client.isTopLevelPath(absolutePath) {
		apiSecret, err = client.topLevelRead(normalizedVaultPath(absolutePath))
	} else {
		apiSecret, err = client.lowLevelRead(normalizedVaultPath(absolutePath))
	}
	if apiSecret != nil {
		secret = NewSecret(apiSecret, absolutePath)
	}
	return secret, err
}

// BatchRead returns secrets for given paths
func (client *Client) BatchRead(absolutePaths []string) (secrets []*Secret, err error) {
	return client.BatchOperation(absolutePaths, OP_READ, make([]*Secret, 0))
}

// Write writes secret to given path, using given Client
func (client *Client) Write(absolutePath string, secret *Secret) (err error) {
	if client.isTopLevelPath(absolutePath) {
		err = client.topLevelWrite(normalizedVaultPath(absolutePath))
	} else {
		err = client.lowLevelWrite(normalizedVaultPath(absolutePath), secret.GetAPISecret())
	}

	return err
}

// BatchWrite writes provided secrets to Vault
func (client *Client) BatchWrite(absolutePaths []string, secrets []*Secret) (err error) {
	_, err = client.BatchOperation(absolutePaths, OP_WRITE, secrets)
	return err
}

// Delete deletes secret at given absolutePath, using given client
func (client *Client) Delete(absolutePath string) (err error) {
	if client.isTopLevelPath(absolutePath) {
		err = client.topLevelDelete(normalizedVaultPath(absolutePath))
	} else {
		err = client.lowLevelDelete(normalizedVaultPath(absolutePath))
	}

	return err
}

// List elements at the given absolutePath, using the given client
func (client *Client) List(absolutePath string) (result []string, err error) {
	defer client.cacheMutex.Unlock()

	client.cacheMutex.Lock()
	if val, ok := client.listCache[absolutePath]; ok {
		return val, nil
	}
	client.cacheMutex.Unlock() // reading from Vault will be relatively slow
	if client.isTopLevelPath(absolutePath) {
		result = client.listTopLevel()
	} else {
		result, err = client.listLowLevel(normalizedVaultPath(absolutePath))
	}
	client.cacheMutex.Lock()
	client.listCache[absolutePath] = result

	return result, err
}

// GetType returns the file type the given absolutePath points to. Possible return values are BACKEND, NODE, LEAF or NONE
func (client *Client) GetType(absolutePath string) (kind PathKind) {
	if client.isTopLevelPath(absolutePath) {
		kind = client.topLevelType(normalizedVaultPath(absolutePath))
	} else {
		kind = client.lowLevelType(normalizedVaultPath(absolutePath))
	}

	return kind
}

// Traverse traverses given absolutePath via DFS and returns sub-paths in array
func (client *Client) Traverse(absolutePath string) (paths []string) {
	if client.isTopLevelPath(absolutePath) {
		paths = client.topLevelTraverse()
	} else {
		paths = client.lowLevelTraverse(normalizedVaultPath(absolutePath))
	}

	return paths
}

// SubpathsForPath will return an array of absolute paths at or below path
func (client *Client) SubpathsForPath(path string) (filePaths []string, err error) {
	switch t := client.GetType(path); t {
	case LEAF:
		filePaths = append(filePaths, filepath.Clean(path))
	case NODE:
		for _, traversedPath := range client.Traverse(path) {
			filePaths = append(filePaths, traversedPath)
		}
	default:
		return filePaths, fmt.Errorf("Not a valid path for operation: %s", path)
	}
	return filePaths, nil
}

// ClearCache clears the list cache
func (client *Client) ClearCache() {
	client.cacheMutex.Lock()
	client.listCache = make(map[string][]string)
	client.cacheMutex.Unlock()
}
