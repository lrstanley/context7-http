// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"fmt"

	"github.com/lrstanley/context7-http/internal/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchLibraryDocsInput struct {
	ResourceURI string   `json:"resourceURI" jsonschema:"Library resource URI (e.g. context7://libraries/<project>), retrieved from resolve-library-uri."`
	Topic       string   `json:"topic,omitempty" jsonschema:"Documentation topic to focus search on (e.g. hooks, routing). Concise and specific, 1-10 words. Strongly encouraged if folders are not provided."`
	Tokens      int      `json:"tokens,omitempty" jsonschema:"Maximum number of tokens of documentation to retrieve. Higher values provide more context but consume more tokens."`
	Folders     []string `json:"folders,omitempty" jsonschema:"List of folders to focus documentation on."`
}

func (s *Server) toolSearchLibraryDocs() *Tool[searchLibraryDocsInput, any] {
	return &Tool[searchLibraryDocsInput, any]{
		Tool: &mcp.Tool{
			Name: "search-library-docs",
			Description: fmt.Sprintf(
				"%s Default token limit: %d.",
				s.mustRender("search_library_docs_desc", nil),
				api.DefaultMinimumDocTokens,
			),
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				IdempotentHint:  false,
				DestructiveHint: new(false),
				OpenWorldHint:   new(true),
			},
		},
		Handler: func(ctx context.Context, _ *mcp.CallToolRequest, input searchLibraryDocsInput) (*mcp.CallToolResult, any, error) {
			result, err := s.client.SearchLibraryDocsText(ctx, input.ResourceURI, &api.SearchLibraryDocsParams{
				Topic:   input.Topic,
				Tokens:  input.Tokens,
				Folders: input.Folders,
			})
			if err != nil {
				s.logger.ErrorContext(ctx, "failed to retrieve library documentation text from Context7", "error", err)
				return nil, "Failed to retrieve library documentation text from Context7.", nil
			}
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: result}}}, nil, nil
		},
	}
}
