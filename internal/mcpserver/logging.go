// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func loggingHooks(logger *slog.Logger, hooks *server.Hooks) *server.Hooks {
	if hooks == nil {
		hooks = &server.Hooks{}
	}

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		l := logger.With(
			slog.Any("id", id),
			slog.String("method", string(method)),
			slog.String("op", "before-any"),
		)

		switch m := message.(type) {
		case *mcp.ReadResourceRequest:
			l = l.With(slog.String("resource", m.Params.URI))
		case *mcp.CallToolRequest:
			l = l.With(slog.String("tool", m.Params.Name))
		default:
			l = l.With(slog.String("type", fmt.Sprintf("%T", m)))
		}

		l.DebugContext(ctx, "received event")
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, _, _ any) {
		logger.DebugContext(
			ctx, "received event",
			slog.Any("id", id),
			slog.String("method", string(method)),
			slog.String("op", "success"),
		)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, _ any, err error) {
		logger.ErrorContext(
			ctx, "error occurred",
			slog.Any("id", id),
			slog.String("method", string(method)),
			slog.String("op", "error"),
			slog.Any("error", err),
		)
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		logger.DebugContext(
			ctx, "received event",
			slog.Any("id", id),
			slog.String("method", message.Method),
			slog.String("op", "before-initialize"),
			slog.String("proto", message.Params.ProtocolVersion),
			slog.String("clientName", message.Params.ClientInfo.Name),
			slog.String("clientVersion", message.Params.ClientInfo.Version),
		)
	})

	return hooks
}
