package client

import (
	"errors"

	"github.com/fishi0x01/vsh/log"
	"github.com/hashicorp/vault/api"
)

func (client *Client) topLevelRead(path string) (secret *api.Secret, err error) {
	err = errors.New(path + " is a backend and cannot be read")
	return secret, err
}

func (client *Client) lowLevelRead(path string) (secret *api.Secret, err error) {
	secret, err = client.Vault.Logical().Read(client.getKVDataPath(path))
	if err != nil {
		log.AppTrace("%+v", err)
	}
	return secret, err
}
