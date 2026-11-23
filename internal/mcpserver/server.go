// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/lrstanley/context7-http/internal/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const instructions = `Context7 provides up-to-date documentation and code examples for libraries and frameworks.

1. Call 'resolve-library-uri' with a library or package name to find its Context7 resource URI.
2. Call 'search-library-docs' with that resource URI (and optionally a topic) to fetch documentation.

Prefer this over prior knowledge when answering questions about a library's API, configuration, or usage, as it reflects the latest published docs.`

type Server struct {
	*mcp.Server

	client        *api.Client
	logger        *slog.Logger
	baseVariables map[string]any
}

func New(_ context.Context, logger *slog.Logger, version string, client *api.Client) (*Server, error) {
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
				Instructions: instructions,
				Logger:       sdkLogger(logger),
				// Tools are static, and we don't send client-side log messages
				// (deprecated as of 2026-07-28), so advertise only what we do.
				Capabilities: &mcp.ServerCapabilities{
					Tools: &mcp.ToolCapabilities{ListChanged: false},
				},
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

// Handler returns the streamable HTTP handler for the server. It runs in
// stateless mode, which is required to serve protocol version 2026-07-28 and
// later, and suits this server as no tool makes server-to-client requests.
// Clients on older protocol versions continue to work, with each request
// handled by a temporary session.
//
// Cross-origin browser requests are rejected, matching the SDK's default
// before v1.6.0.
func (s *Server) Handler() http.Handler {
	return http.NewCrossOriginProtection().Handler(
		mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
			return s.Server
		}, &mcp.StreamableHTTPOptions{
			Stateless:                    true,
			Logger:                       sdkLogger(s.logger),
			PropagateRequestCancellation: true,
		}),
	)
}

// toolError converts an error from the API client into one that is safe to
// return to the client as a tool execution error. Errors the caller can act on
// are passed through, and anything else is logged and replaced with a generic
// message.
func (s *Server) toolError(ctx context.Context, msg string, err error) error {
	if errors.Is(err, api.ErrRateLimited) || errors.Is(err, api.ErrInvalidResourceURI) {
		return err
	}
	s.logger.ErrorContext(ctx, "tool call failed", "msg", msg, "error", err)
	return errors.New(msg)
}
