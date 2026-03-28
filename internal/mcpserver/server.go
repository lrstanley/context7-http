// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/lrstanley/context7-http/internal/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	*mcp.Server
	client        *api.Client
	logger        *slog.Logger
	baseVariables map[string]any
}

func New(_ context.Context, logger *slog.Logger, version string, keepAlive time.Duration, client *api.Client) (*Server, error) {
	name := "Context7"
	srv := &Server{
		client: client,
		logger: logger,
		Server: mcp.NewServer(
			&mcp.Implementation{
				Name:       name,
				Title:      "Context7",
				WebsiteURL: "https://context7.com",
				Version:    version,
			}, &mcp.ServerOptions{
				Instructions: "TODO",
				Logger:       logger,
				KeepAlive:    keepAlive,
			},
		),
	}

	srv.AddReceivingMiddleware(loggingMiddleware(logger))

	// Tools.
	srv.toolResolveLibraryID().Add(srv.Server)
	srv.toolSearchLibraryDocs().Add(srv.Server)

	srv.baseVariables = map[string]any{
		"ServerName":    name,
		"ServerVersion": version,
	}

	return srv, nil
}

func (s *Server) Handler() http.Handler {
	return mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return s.Server
	}, nil)
}
