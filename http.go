// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/context7-http/internal/mcpserver"
	"github.com/mark3labs/mcp-go/server"
)

func httpServer(logger *slog.Logger, srv *mcpserver.Server) *http.Server {
	r := chi.NewRouter()

	r.Use(
		chix.NewConfig().
			SetLogger(logger).
			Use(),
	)

	if len(cli.Flags.TrustedProxies) > 0 {
		r.Use(chix.UseRealIPStringOpts(cli.Flags.TrustedProxies))
	}

	// Core middeware.
	r.Use(
		chix.UseContextIP(),
		chix.UseStructuredLogger(chix.DefaultLogConfig()),
		chix.UseStripSlashes(),
		middleware.Compress(5),
		chix.UseSecurityText(chix.SecurityTextConfig{
			ExpiresIn: 182 * 24 * time.Hour,
			Contacts: []string{
				"mailto:liam@liam.sh",
				"https://liam.sh/chat",
				"https://github.com/lrstanley",
			},
			KeyLinks:  []string{"https://github.com/lrstanley.gpg"},
			Languages: []string{"en"},
		}),
	)

	sseServer := server.NewSSEServer(
		srv.MCPServer,
		server.WithBaseURL(cli.Flags.BaseURL),
	)
	r.Handle("/sse", sseServer)
	r.Handle("/message", sseServer)

	streamableServer := server.NewStreamableHTTPServer(
		srv.MCPServer,
		server.WithHeartbeatInterval(cli.Flags.HeartbeatInterval),
	)
	r.Handle("/mcp", streamableServer)

	if cli.Debug {
		r.With(chix.UsePrivateIP()).Mount("/debug", middleware.Profiler())
	}

	r.With(middleware.ThrottleBacklog(1, 5, 5*time.Second)).Get("/healthy", func(w http.ResponseWriter, r *http.Request) {
		chix.JSON(w, r, 200, chix.M{
			"status": "ok",
		})
	})

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body style="background-color:#383838;"><h1 style="color:white;">Context7 MCP Server</h1><ul>`))
		for _, link := range cli.GetVersion().AppInfo.Links {
			_, _ = fmt.Fprintf(w, `<li><a style="color:white;text-transform:capitalize;" href=%q>%s</a></li>`, link.URL, link.Name)
		}
		_, _ = fmt.Fprintf(w, `<li><a style="color:white;" href=%q>SSE -- <code>%s/sse</code></a></li>`, cli.Flags.BaseURL+"/sse", cli.Flags.BaseURL)
		_, _ = fmt.Fprintf(w, `<li><a style="color:white;" href=%q>MCP -- <code>%s/mcp</code></a></li>`, cli.Flags.BaseURL+"/mcp", cli.Flags.BaseURL)
		_, _ = w.Write([]byte(`</ul></body></html>`))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		chix.ErrorWithCode(w, r, http.StatusNotFound)
	})

	return &http.Server{
		Addr:    cli.Flags.BindAddr,
		Handler: r,
		// Must explicitly stay set to 0 because long-lived connections.
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
}
