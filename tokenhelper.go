//go:build !notokenhelper
// +build !notokenhelper

package main

import (
	"github.com/hashicorp/vault/api/cliconfig"
)

func getTokenFromHelper() (string, error) {
	helper, err := cliconfig.DefaultTokenHelper()
	if err != nil {
		return "", err
	}
	return helper.Get()
}
