package client

import (
	"strings"

	"github.com/fishi0x01/vsh/log"
)

func (client *Client) topLevelTraverse() (result []string) {
	for k := range client.KVBackends {
		result = append(result, k)
	}

	return result
}

func (client *Client) lowLevelTraverse(path string, shallow bool) (result []string) {
	s, err := client.cache.List(client.getKVMetaDataPath(path))
	if err != nil {
		log.AppTrace("%+v", err)
		return
	}

	if s != nil {
		if keysInterface, ok := s.Data["keys"]; ok {
			for _, valInterface := range keysInterface.([]interface{}) {
				val := valInterface.(string)
				// prevent ambiguous dir/file to be added twice
				if strings.HasSuffix(val, "/") {
					// dir
					if !shallow {
						result = append(result, client.lowLevelTraverse(path+"/"+val, false)...)
					}
				} else {
					// file
					leaf := strings.ReplaceAll("/"+path+"/"+val, "//", "/")
					result = append(result, leaf)
				}
			}
		}
	} else {
		leaf := strings.ReplaceAll("/"+path, "//", "/")
		result = append(result, leaf)
	}
	return result
}
