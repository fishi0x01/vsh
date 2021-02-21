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

func (client *Client) isAmbiguous(path string, dirFiles map[string]int) (result bool) {
	// check if path exists as file and directory
	result = false
	if _, ok := dirFiles[filepath.Base(path)]; ok {
		if _, ok := dirFiles[filepath.Base(path)+"/"]; ok {
			result = true
		}
	}
	return result
}

func (client *Client) getDirFiles(path string) (result map[string]int) {
	// get current directory content
	result = make(map[string]int)
	pathTrim := strings.TrimSuffix(path, "/")
	lsPath := client.getKVMetaDataPath(filepath.Dir(pathTrim))
	s, err := client.cache.List(lsPath)
	if err == nil && s != nil {
		if keysInterface, ok := s.Data["keys"]; ok {
			for _, valInterface := range keysInterface.([]interface{}) {
				val := valInterface.(string)
				result[val] = 1
			}
		}
	}
	return result
}

func (client *Client) lowLevelType(path string) (result PathKind) {
	dirFiles := client.getDirFiles(path)
	if client.isAmbiguous(path, dirFiles) {
		if strings.HasSuffix(path, "/") {
			result = NODE
		} else {
			result = LEAF
		}
	} else {
		hasNode := false
		kvPath := client.getKVMetaDataPath(path + "/")
		s, err := client.cache.List(kvPath)
		if err == nil && s != nil {
			if _, ok := s.Data["keys"]; ok {
				hasNode = true
			}
		}
		if hasNode {
			result = NODE
		} else {
			if _, ok := dirFiles[filepath.Base(path)]; ok {
				result = LEAF
			} else {
				result = NONE
			}
		}
	}
	return result
}
