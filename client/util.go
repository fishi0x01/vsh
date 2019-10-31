package client

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
)

func (client *Client) discoverMountBackend(backend string) {
	if _, ok := client.KVBackends[backend]; !ok {
		// backend is unknown - check if exists
		_, err := client.Vault.Help(backend)
		if err != nil {
			return
		}

		// backend exists
		s, err := client.Vault.Logical().List(backend)
		if err != nil {
			return
		}

		if s == nil {
			return
		}

		if s.Warnings != nil {
			client.KVBackends[backend] = 2
		} else {
			client.KVBackends[backend] = 1
		}
	}
}

func (client *Client) getKVVersion(path string) int {
	mntPath := ""
	if strings.HasPrefix(path, "/") {
		mntPath = strings.Split(path, "/")[1] + "/"
	} else {
		mntPath = strings.Split(path, "/")[0] + "/"
	}

	client.discoverMountBackend(mntPath)

	if version, ok := client.KVBackends[mntPath]; ok {
		return version
	}
	return -1
}

func (client *Client) kvPath(path string, prefix string) string {
	v := client.getKVVersion(path)
	switch v {
	case 1:
		return path
	case 2:
		// https://www.vaultproject.io/docs/secrets/kv/kv-v2.html#acl-rules
		s := strings.SplitN(path, "/", 2)
		if len(s) != 2 {
			panic(fmt.Errorf("Could not properly split path '%s'", path))
		}
		return s[0] + prefix + s[1]
	default:
		return ""
	}
}

func (client *Client) getKVMetaDataPath(path string) string {
	return client.kvPath(path, "/metadata/")
}

func (client *Client) getKVDataPath(path string) string {
	return client.kvPath(path, "/data/")
}

func (client *Client) isTopLevelPath(absolutePath string) bool {
	if strings.Count(absolutePath, "/") < 2 {
		return true
	}
	return false
}

func isValidKV2Data(secret *api.Secret) bool {
	_, exists := secret.Data["data"]
	return exists
}

func transformToKV1Secret(secret api.Secret) *api.Secret {
	secret.Data = secret.Data["data"].(map[string]interface{})
	return &secret
}

func transformToKV2Secret(secret api.Secret) *api.Secret {
	secret.Data = map[string]interface{}{
		"data": secret.Data,
	}

	return &secret
}

func normalizedVaultPath(absolutePath string) string {
	// remove trailing '/'
	return absolutePath[1:]
}

func sliceContains(arr []string, search string) bool {
	for _, s := range arr {
		if s == search {
			return true
		}
	}
	return false
}
