// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"os"
	"time"

	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/context7-http/internal/api"
	"github.com/lrstanley/context7-http/internal/mcpserver"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"

	cli = clix.NewWithDefaults(
		clix.WithAppInfo[Flags](clix.AppInfo{
			Version: version,
			Commit:  commit,
			Date:    date,
			Links:   clix.GithubLinks("github.com/lrstanley/context7-http", "master", "https://liam.sh"),
		}),
	)
)

type Flags struct {
	BindAddr          string        `name:"bind-addr"          env:"BIND_ADDR" default:":8080"`
	BaseURL           string        `name:"base-url"           env:"BASE_URL" default:"http://localhost:8080"`
	TrustedProxies    []string      `name:"trusted-proxies"    env:"TRUSTED_PROXIES" help:"CIDR ranges that we trust the X-Forwarded-For header from"`
	HeartbeatInterval time.Duration `name:"heartbeat-interval" env:"HEARTBEAT_INTERVAL"`
}

func main() {
	logger := cli.GetLogger()
	ctx := context.Background()

	var client *api.Client
	var srv *mcpserver.Server
	var err error

	client, err = api.New(ctx, logger, nil)
	if err != nil {
		logger.ErrorContext(ctx, "failed to initialize api client", "error", err)
		os.Exit(1)
	}

	srv, err = mcpserver.New(ctx, logger, version, client)
	if err != nil {
		logger.ErrorContext(ctx, "failed to initialize mcp server", "error", err)
		os.Exit(1)
	}

	err = chix.NewServerWithoutDefaults(ctx, logger, httpServer(logger, srv), "", "")
	if err != nil {
		logger.ErrorContext(ctx, "http server invocation failed", "error", err)
		os.Exit(1)
	}
}
