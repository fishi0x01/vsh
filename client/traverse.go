package client

import (
	"github.com/fishi0x01/vsh/log"
	"strings"
)

func (client *Client) topLevelTraverse(c chan<- string) {
	for k := range client.KVBackends {
		c <- k
	}
}

func (client *Client) lowLevelTraverse(path string, c chan<- string) {
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
					client.lowLevelTraverse(path+"/"+val, c)
				} else {
					// file
					leaf := strings.ReplaceAll("/"+path+"/"+val, "//", "/")
					c <- leaf
				}
			}
		}
	} else {
		leaf := strings.ReplaceAll("/"+path, "//", "/")
		c <- leaf
	}
}
