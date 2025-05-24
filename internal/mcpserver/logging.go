// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func loggingHooks(hooks *server.Hooks) *server.Hooks {
	if hooks == nil {
		hooks = &server.Hooks{}
	}

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		fields := log.Fields{
			"id":     id,
			"method": method,
			"op":     "before-any",
		}

		switch m := message.(type) {
		case *mcp.ReadResourceRequest:
			fields["resource"] = m.Params.URI
		case *mcp.CallToolRequest:
			fields["tool"] = m.Params.Name
		default:
			fields["type"] = fmt.Sprintf("%T", m)
		}

		// fmt.Printf("%#v\n", message)
		log.FromContext(ctx).WithFields(fields).Debug("received event")
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		log.FromContext(ctx).WithFields(log.Fields{
			"id":     id,
			"method": method,
			"op":     "success",
		}).Debug("received event")
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		log.FromContext(ctx).WithFields(log.Fields{
			"id":     id,
			"method": method,
			"op":     "error",
		}).WithError(err).Error("error occurred")
	})
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		log.FromContext(ctx).WithFields(log.Fields{
			"id":            id,
			"method":        message.Method,
			"op":            "before-initialize",
			"proto":         message.Params.ProtocolVersion,
			"clientName":    message.Params.ClientInfo.Name,
			"clientVersion": message.Params.ClientInfo.Version,
		}).Debug("received event")
	})

	return hooks
}
