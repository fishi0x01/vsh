package client

import (
	"errors"
	"github.com/hashicorp/vault/api"
)

func (client *Client) topLevelWrite(path string) error {
	return errors.New(path + " is a backend and cannot be written")
}

func (client *Client) lowLevelWrite(path string, secret *api.Secret) (err error) {
	if client.getKVVersion(path) == 1 {
		if isValidKV2Data(secret) {
			secret = transformToKV1Secret(*secret)
		}
	}

	if client.getKVVersion(path) == 2 {
		if !isValidKV2Data(secret) {
			secret = transformToKV2Secret(*secret)
		}
	}

	_, err = client.Vault.Logical().Write(client.getKVDataPath(path), secret.Data)
	return err
}
