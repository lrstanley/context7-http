// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"

	"github.com/lrstanley/context7-http/internal/api"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	*server.MCPServer
	client        *api.Client
	baseVariables map[string]any
}

func New(_ context.Context, version string, client *api.Client) (*Server, error) {
	name := "Context7"
	srv := &Server{
		client: client,
		MCPServer: server.NewMCPServer(
			name,
			version,
			// server.WithInstructions(),
			server.WithRecovery(),
			server.WithToolCapabilities(false),
			server.WithHooks(loggingHooks(nil)),
			server.WithPaginationLimit(250),
		),
	}

	srv.AddTool(srv.toolResolveLibraryID())
	srv.AddTool(srv.toolSearchLibraryDocs())
	srv.AddResource(srv.resourceLibrariesAll())
	srv.AddResource(srv.resourceLibrariesTop(500))
	srv.AddResource(srv.resourceLibrariesTop(1000))
	srv.AddResourceTemplate(srv.resourceLibrary()) // TODO: is this just for searching?

	srv.baseVariables = map[string]any{
		"ServerName":    name,
		"ServerVersion": version,
	}

	return srv, nil
}
