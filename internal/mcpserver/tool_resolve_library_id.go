// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type resolveLibraryIDInput struct {
	LibraryName string `json:"libraryName" jsonschema:"Library name to search for, returning a context7-compatible library resource URI."`
}

func (s *Server) toolResolveLibraryID() *Tool[resolveLibraryIDInput, any] {
	return &Tool[resolveLibraryIDInput, any]{
		Tool: &mcp.Tool{
			Name:        "resolve-library-uri",
			Description: s.mustRender("resolve_library_id_desc", nil),
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				IdempotentHint:  false,
				DestructiveHint: new(false),
				OpenWorldHint:   new(false),
			},
		},
		Handler: func(ctx context.Context, _ *mcp.CallToolRequest, input resolveLibraryIDInput) (*mcp.CallToolResult, any, error) {
			results, err := s.client.SearchLibraries(ctx, input.LibraryName)
			if err != nil {
				s.logger.ErrorContext(ctx, "failed to retrieve library documentation data from Context7", "error", err)
				return nil, "Failed to retrieve library documentation data from Context7.", nil
			}

			if len(results) == 0 {
				return nil, "No documentation libraries available matching that criteria.", nil
			}

			content := make([]mcp.Content, 0, len(results))
			var text string

			for _, result := range results {
				text, err = s.render("resolve_library_id_resp", result)
				if err != nil {
					return nil, "", err
				}
				content = append(content, &mcp.TextContent{Text: text})
			}
			return &mcp.CallToolResult{Content: content}, nil, nil
		},
	}
}
