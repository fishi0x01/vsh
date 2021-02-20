package client

import (
	"path/filepath"
	"strings"
)

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

func (client *Client) FilterPaths(paths []string, kind PathKind) (filtered []string) {
	for _, path := range paths {
		if client.GetType(path) == kind {
			filtered = append(filtered, path)
		}
	}
	return filtered
}

var cachedPath = ""
var cachedDirFiles = make(map[string]int)

func (client *Client) isAmbiguous(path string) (result bool) {
	// get current directory content
	if cachedPath != path {
		pathTrim := strings.TrimSuffix(path, "/")
		cachedDirFiles = make(map[string]int)
		s, err := client.Vault.Logical().List(client.getKVMetaDataPath(filepath.Dir(pathTrim)))
		if err == nil && s != nil {
			if keysInterface, ok := s.Data["keys"]; ok {
				for _, valInterface := range keysInterface.([]interface{}) {
					val := valInterface.(string)
					cachedDirFiles[val] = 1
				}
			}
		}
		cachedPath = path
	}

	// check if path exists as file and directory
	result = false
	if _, ok := cachedDirFiles[filepath.Base(path)]; ok {
		if _, ok := cachedDirFiles[filepath.Base(path)+"/"]; ok {
			result = true
		}
	}
	return result
}

func (client *Client) lowLevelType(path string) (result PathKind) {
	if client.isAmbiguous(path) {
		if strings.HasSuffix(path, "/") {
			result = NODE
		} else {
			result = LEAF
		}
	} else {
		hasNode := false
		s, err := client.Vault.Logical().List(client.getKVMetaDataPath(path + "/"))
		if err == nil && s != nil {
			if _, ok := s.Data["keys"]; ok {
				hasNode = true
			}
		}
		if hasNode {
			result = NODE
		} else {
			if _, ok := cachedDirFiles[filepath.Base(path)]; ok {
				result = LEAF
			} else {
				result = NONE
			}
		}
	}
	return result
}
