package client

import (
	"errors"
	"github.com/fishi0x01/vsh/log"
)

func (client *Client) topLevelDelete(path string) error {
	return errors.New(path + " is a vault backend and cannot be deleted")
}

func (client *Client) lowLevelDelete(path string) (err error) {
	_, err = client.Vault.Logical().Delete(client.getKVMetaDataPath(path))
	if err != nil {
		log.AppTrace("%+v", err)
	}
	return err
}
