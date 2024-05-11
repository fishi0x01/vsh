// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugincontainer

import (
	"context"
	"fmt"
	"runtime"

	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin/runner"
)

func ReattachFunc(logger hclog.Logger, id, hostSocketDir string) (runner.AttachedRunner, error) {
	if runtime.GOOS != "linux" {
		return nil, errUnsupportedOS
	}

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	_, err = client.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("container with ID %s not found: %w", id, err)
	}

	return &containerRunner{
		dockerClient:  client,
		logger:        logger,
		id:            id,
		hostSocketDir: hostSocketDir,
	}, nil
}
