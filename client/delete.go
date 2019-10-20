package client

import (
	"errors"
)

func (client *Client) topLevelDelete(path string) error {
	return errors.New(path + " is a vault backend and cannot be deleted")
}

func (client *Client) lowLevelDelete(path string) (err error) {
	_, err = client.Vault.Logical().Delete(client.getKVMetaDataPath(path))
	return err
}
