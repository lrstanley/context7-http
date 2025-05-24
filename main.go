// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/context7-http/internal/api"
	"github.com/lrstanley/context7-http/internal/mcpserver"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"

	cli = &clix.CLI[Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/context7-http", "master", "https://liam.sh"),
		VersionInfo: &clix.VersionInfo[Flags]{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
	}

	srv    *mcpserver.Server
	client *api.Client
)

type Flags struct {
	BindAddr          string        `long:"bind-addr"          env:"BIND_ADDR"          default:":8080"`
	BaseURL           string        `long:"base-url"           env:"BASE_URL"           default:"http://localhost:8080"`
	TrustedProxies    []string      `long:"trusted-proxies"    env:"TRUSTED_PROXIES"    description:"CIDR ranges that we trust the X-Forwarded-For header from"`
	HeartbeatInterval time.Duration `long:"heartbeat-interval" env:"HEARTBEAT_INTERVAL"`
}

func main() {
	cli.LoggerConfig.Pretty = true
	cli.Parse()

	ctx := context.Background()

	var err error

	client, err = api.New(ctx, nil)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize api client")
	}

	srv, err = mcpserver.New(ctx, version, client)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize mcp server")
	}

	chix.SetServerDefaults = false
	if err = chix.RunContext(ctx, httpServer(ctx)); err != nil {
		log.FromContext(ctx).WithError(err).Fatal("shutting down")
	}
}
