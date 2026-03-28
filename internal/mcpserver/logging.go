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

func loggingMiddleware(logger *slog.Logger) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			l := logger.With(
				slog.String("session", req.GetSession().ID()),
				slog.String("method", method),
				slog.String("type", fmt.Sprintf("%T", req)),
			)

			switch r := any(req).(type) {
			case *mcp.CallToolParamsRaw:
				l.With(slog.String("tool_name", r.Name))
			case *mcp.ReadResourceParams:
				l.With(slog.String("resource_uri", r.URI))
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
