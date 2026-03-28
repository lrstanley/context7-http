// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool[I any, O any] struct {
	Tool    *mcp.Tool
	Handler mcp.ToolHandlerFor[I, O]
}

func (t *Tool[I, O]) Add(server *mcp.Server) {
	mcp.AddTool(server, t.Tool, t.Handler)
}
