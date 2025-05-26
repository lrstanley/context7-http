// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"

	"github.com/apex/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ResolveLibraryIDParams struct {
	LibraryName string `json:"libraryName"`
}

func (s *Server) toolResolveLibraryID() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"resolve-library-uri",
		mcp.WithString(
			"libraryName",
			mcp.Required(),
			mcp.MinLength(2),
			mcp.MaxLength(100),
			mcp.Description("Library name to search for, returning a context7-compatible library resource URI."),
		),
		mcp.WithDescription(s.mustRender("resolve_library_id_desc", nil)),
	)

	return tool, mcp.NewTypedToolHandler(func(ctx context.Context, _ mcp.CallToolRequest, params ResolveLibraryIDParams) (*mcp.CallToolResult, error) {
		results, err := s.client.SearchLibraries(ctx, params.LibraryName)
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("failed to retrieve library documentation data from Context7")
			return mcp.NewToolResultError("Failed to retrieve library documentation data from Context7."), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultError("No documentation libraries available matching that criteria."), nil
		}
		resp, err := s.render("resolve_library_id_resp", map[string]any{
			"Results": results,
		})
		return mcp.NewToolResultText(resp), err
	})
}
