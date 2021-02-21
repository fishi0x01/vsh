package client

import (
	"errors"
	"github.com/fishi0x01/vsh/log"
)

func (client *Client) listTopLevel() (result []string) {
	for k := range client.KVBackends {
		result = append(result, k)
	}
	return result
}

func (client *Client) listLowLevel(path string) (result []string, err error) {
	t := client.lowLevelType(path)
	if t != BACKEND && t != NODE {
		return nil, errors.New("Not a directory: " + path)
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
