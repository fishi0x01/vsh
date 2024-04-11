//go:build !notokenhelper
// +build !notokenhelper

package main

import (
	"github.com/hashicorp/vault/command/config"
)

func getTokenFromHelper() (string, error) {
	helper, err := config.DefaultTokenHelper()
	if err != nil {
		return "", err
	}
	return helper.Get()
}
