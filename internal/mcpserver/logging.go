// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// sdkLogger returns a logger for the MCP SDK. In stateless mode the SDK logs
// session connect/disconnect at INFO on every request, so below DEBUG only
// warnings and errors are passed through.
func sdkLogger(logger *slog.Logger) *slog.Logger {
	if logger.Enabled(context.Background(), slog.LevelDebug) {
		return logger
	}
	return slog.New(minLevelHandler{Handler: logger.Handler(), level: slog.LevelWarn})
}

type minLevelHandler struct {
	slog.Handler

	level slog.Level
}

func (h minLevelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level && h.Handler.Enabled(ctx, level)
}

func (h minLevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return minLevelHandler{Handler: h.Handler.WithAttrs(attrs), level: h.level}
}

func (h minLevelHandler) WithGroup(name string) slog.Handler {
	return minLevelHandler{Handler: h.Handler.WithGroup(name), level: h.level}
}

// serverRequest is implemented by all [mcp.ServerRequest] types.
type serverRequest interface {
	ProtocolVersion() string
	ClientInfo() *mcp.Implementation
}

func loggingMiddleware(logger *slog.Logger) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			l := logger.With(
				slog.String("method", method),
				slog.String("type", fmt.Sprintf("%T", req)),
			)

			// Sessions are temporary in stateless mode, so the session ID is
			// usually empty. Log the per-request client metadata instead.
			if r, ok := req.(serverRequest); ok {
				if v := r.ProtocolVersion(); v != "" {
					l = l.With(slog.String("protocol_version", v))
				}
				if info := r.ClientInfo(); info != nil {
					l = l.With(
						slog.String("client_name", info.Name),
						slog.String("client_version", info.Version),
					)
				}
			}

			switch p := req.GetParams().(type) {
			case *mcp.CallToolParamsRaw:
				l = l.With(slog.String("tool_name", p.Name))
			case *mcp.ReadResourceParams:
				l = l.With(slog.String("resource_uri", p.URI))
			}

			start := time.Now()

			l.DebugContext(ctx, "received event")

			result, err := next(ctx, method, req)
			l = l.With(slog.Duration("duration", time.Since(start)))

			if err != nil {
				l.ErrorContext(
					ctx, "error occurred",
					slog.Any("error", err),
				)
			} else {
				l.DebugContext(ctx, "completed event")
			}

			return result, err
		}
	}
}
