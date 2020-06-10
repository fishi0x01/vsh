package client

import "strings"

// PathKind describes the type of a path
type PathKind int

// types of paths
const (
	BACKEND PathKind = iota
	NODE
	LEAF
	NONE
)

func (client *Client) topLevelType(path string) PathKind {
	if path == "" {
		return BACKEND
	} else if _, ok := client.KVBackends[path+"/"]; ok {
		return BACKEND
	} else {
		return NONE
	}
}

func (client *Client) lowLevelType(path string) PathKind {
	is_node := false
	is_leaf := false

	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err != nil {
		return NONE
	}

	if s != nil {
		is_node = true
	}

	s, err = client.Vault.Logical().Read(client.getKVDataPath(path))
	if err == nil && s != nil {
		is_leaf = true
	}

	if is_leaf && !is_node {
		return LEAF
	}

	if is_node && !is_leaf {
		return NODE
	}

	if is_leaf && is_node {
		// vault namespace path can overlap with a key, e.g.,
		// secret/a and secret/a/b
		// --> in that case, we have a leaf and a node when checking secret/a
		if strings.HasSuffix(path, "/") {
			return NODE
		} else {
			return LEAF
		}
	}

	return NONE
}
