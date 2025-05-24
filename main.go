// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"

	"github.com/apex/log"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/context7-http/internal/api"
	"github.com/lrstanley/context7-http/internal/mcpserver"
	"github.com/lrstanley/context7-http/internal/session"
	"github.com/mark3labs/mcp-go/server"
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

	client *api.Client
)

type Flags struct {
	BindAddr string `long:"bind-addr" env:"BIND_ADDR" default:":8080"`
	BaseURL  string `long:"base-url" env:"BASE_URL" default:"http://localhost:8080"`
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

	srv, err := mcpserver.New(ctx, version, client)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize mcp server")
	}

	sse := server.NewSSEServer(
		srv.MCPServer,
		server.WithBaseURL(cli.Flags.BaseURL),
		server.WithHTTPContextFunc(session.AddIDToContext),
		server.WithAppendQueryToMessageEndpoint(),
	)

	log.FromContext(ctx).WithField("bind-addr", cli.Flags.BindAddr).Info("starting sse server")
	err = sse.Start(cli.Flags.BindAddr)
	if err != nil {
		log.WithError(err).Fatal("failed running sse server")
	}
}
