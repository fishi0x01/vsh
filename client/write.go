package client

import (
	"errors"
	"github.com/hashicorp/vault/api"
)

func (client *Client) topLevelWrite(path string) error {
	return errors.New(path + " is a backend and cannot be written")
}

func (client *Client) lowLevelWrite(path string, secret *api.Secret) (err error) {
	_, err = client.Vault.Logical().Write(client.getKVDataPath(path), secret.Data)
	return err
}