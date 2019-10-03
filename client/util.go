package client

import (
	"fmt"
	"strings"
)

func (client *Client) getKVVersion(path string) int {
	mntPath := ""
	if strings.HasPrefix(path, "/") {
		mntPath = strings.Split(path, "/")[1] + "/"
	} else {
		mntPath = strings.Split(path, "/")[0] + "/"
	}
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
		if strings.Contains(path, "/") {
			panic(fmt.Errorf("Unknown KV Version '%v' for path '%s'", v, path))
		}
		// we are in the root path
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
