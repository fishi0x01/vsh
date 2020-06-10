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

func (client *Client) lowLevelType(path string) (result PathKind) {
	result = NONE
	isNode := false
	isLeaf := false

	s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path))
	if err == nil && s != nil {
		isNode = true
	}

	s, err = client.Vault.Logical().Read(client.getKVDataPath(path))
	if err == nil && s != nil {
		isLeaf = true
	}

	if isLeaf && !isNode {
		result = LEAF
	}

	if isNode && !isLeaf {
		result = NODE
	}

	if isLeaf && isNode {
		// vault namespace path can overlap with a key, e.g.,
		// secret/a and secret/a/b
		// --> in that case, we have a leaf and a node when checking secret/a
		if strings.HasSuffix(path, "/") {
			result = NODE
		} else {
			result = LEAF
		}
	}

	return
}
