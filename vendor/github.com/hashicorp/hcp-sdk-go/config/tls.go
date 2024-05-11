// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import "crypto/tls"

// cloneTLSConfig will shallow clone a TLS. tls.Config already has a Clone method
// but versions before Go 1.15 did not allow to clone nil.
// TODO: remove once the SDK only supports Go 1.16 and newer.
func cloneTLSConfig(original *tls.Config) *tls.Config {
	if original == nil {
		return nil
	}

	return original.Clone()
}
