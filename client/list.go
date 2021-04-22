package client

import (
	"errors"
	_"fmt"

	"github.com/fishi0x01/vsh/log"
)

const (
	MODE_IGNORE_NONE_DIRECTORY = 0
	MODE_DIRECTORY_ONLY = 1
)

func (client *Client) listTopLevel() (result []string) {
	for k := range client.KVBackends {
		result = append(result, k)
	}
	return result
}

func (client *Client) listLowLevel(path string, mode int) (result []string, err error) {
	t := client.lowLevelType(path)
	if t != BACKEND && t != NODE {
		if mode == MODE_DIRECTORY_ONLY {
			return nil, errors.New("Not a directory: " + path)
		} else if mode == MODE_IGNORE_NONE_DIRECTORY {
			return nil, nil
		}
	}

	s, err := client.cache.List(client.getKVMetaDataPath(path))
	if err != nil {
		log.AppTrace("%+v", err)
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

func (client *Client) listAllLowLevel(path string, subResult* []string) (result[]string, err error) {
	subResult2, err := client.listLowLevel(path, MODE_IGNORE_NONE_DIRECTORY)
	if subResult2 != nil {
		var subResult3 []string
		for _, resultPath := range subResult2 {
			_, err := client.listAllLowLevel(normalizedVaultPath(path + resultPath), subResult)
			if err != nil {
				log.AppTrace("%+v", err)
				return subResult3, err
			}
			result = append(result, normalizedVaultPath(path + resultPath))
		}
	} else {
		*subResult = append(*subResult, path)
	}
	return *subResult, err
}

func (client *Client) listAllFromTopLevel() (result []string, err error) {
	res:= []string{}
	for k := range client.KVBackends {
		result, err = client.listAllLowLevel(normalizedVaultPath(k), &res)
	}
	return result, err
}