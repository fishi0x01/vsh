package client

import (
	"github.com/hashicorp/vault/api"
)

// Secret holds vault secret and offers operations to simplify KV abstraction
type Secret struct {
	vaultSecret *api.Secret
	Path        string
}

// NewSecret create a new Secret object
func NewSecret(vaultSecret *api.Secret, path string) *Secret {
	return &Secret{
		vaultSecret: vaultSecret,
		Path:        path,
	}
}

// GetAPISecret getter method for vault secret in Secret object
func (secret *Secret) GetAPISecret() *api.Secret {
	return secret.vaultSecret
}

// GetData returns the secret data as a map and is KV agnostic
func (secret *Secret) GetData() map[string]interface{} {
	data := make(map[string]interface{})
	for k, v := range secret.vaultSecret.Data {
		if rec, ok := v.(map[string]interface{}); ok {
			// KV 2
			if k == "data" {
				for kk, vv := range rec {
					data[kk] = vv
				}
			}
		} else {
			// KV 1
			data[k] = v
		}
	}
	return data
}

// SetData set given data as vault secret data and is KV agnostic
func (secret *Secret) SetData(data map[string]interface{}) {
	isKV2 := false
	if val, hasData := secret.vaultSecret.Data["data"]; hasData {
		if _, isKV2 := val.(map[string]interface{}); isKV2 {
			// KV2
			secret.vaultSecret.Data["data"] = data
		}
	}
	if !isKV2 {
		// KV1
		secret.vaultSecret.Data = data
	}
}
