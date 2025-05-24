// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/lrstanley/context7-http/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SearchLibraryDocsParams struct {
	ResourceURI string   `json:"resourceURI"`
	Topic       string   `json:"topic,omitempty"`
	Tokens      int      `json:"tokens,omitempty"`
	Folders     []string `json:"folders,omitempty"`
}

func (s *Server) toolSearchLibraryDocs() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	tool = mcp.NewTool(
		"search-library-docs",
		mcp.WithString(
			"resourceURI",
			mcp.Required(),
			mcp.Description("Library resource URI (e.g., 'context7://libraries/<project>'), which is retrieved from 'resolve-library-uri'."),
		),
		mcp.WithString(
			"topic",
			mcp.Description(
				"Documentation topic to focus search on (e.g., 'hooks', 'routing')."+
					" This should be concise and specific, 1-10 words if possible."+
					" This is strongly encouraged to be provided if folders are not provided.",
			),
		),
		mcp.WithNumber(
			"tokens",
			mcp.Description(
				(fmt.Sprintf("Maximum number of tokens of documentation to retrieve (default: %d).", api.DefaultMinimumDocTokens))+
					" Higher values provide more context but consume more tokens.",
			),
		),
		mcp.WithArray(
			"folders",
			mcp.Description("List of folders to focus documentation on."),
		),
		mcp.WithDescription(s.mustRender("search_library_docs_desc", nil)),
	)

	return tool, mcp.NewTypedToolHandler(func(ctx context.Context, request mcp.CallToolRequest, params SearchLibraryDocsParams) (*mcp.CallToolResult, error) {
		result, err := s.client.SearchLibraryDocsText(ctx, params.ResourceURI, &api.SearchLibraryDocsParams{
			Topic:   params.Topic,
			Tokens:  params.Tokens,
			Folders: params.Folders,
		})
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("failed to retrieve library documentation text from Context7")
			return mcp.NewToolResultError("Failed to retrieve library documentation text from Context7."), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}
