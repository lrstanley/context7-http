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

type resolveLibraryIDOutput struct {
	Libraries []resolveLibraryIDResult `json:"libraries" jsonschema:"Libraries matching the query, sorted by relevance."`
}

type resolveLibraryIDResult struct {
	ResourceURI   string  `json:"resourceURI" jsonschema:"Context7-compatible library resource URI, for use with search-library-docs."`
	Title         string  `json:"title" jsonschema:"Library or package name."`
	Description   string  `json:"description,omitempty" jsonschema:"Short summary of the library."`
	TotalSnippets int     `json:"totalSnippets,omitempty" jsonschema:"Number of available code examples."`
	TrustScore    float64 `json:"trustScore,omitempty" jsonschema:"Authority indicator, from 0 to 10."`
}

func (s *Server) toolResolveLibraryID() *Tool[resolveLibraryIDInput, *resolveLibraryIDOutput] {
	return &Tool[resolveLibraryIDInput, *resolveLibraryIDOutput]{
		Tool: &mcp.Tool{
			Name:        "resolve-library-uri",
			Title:       "Resolve Library URI",
			Description: s.mustRender("resolve_library_id_desc", nil),
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				IdempotentHint:  false,
				DestructiveHint: new(false),
				OpenWorldHint:   new(false),
			},
		},
		Handler: func(ctx context.Context, _ *mcp.CallToolRequest, input resolveLibraryIDInput) (*mcp.CallToolResult, *resolveLibraryIDOutput, error) {
			results, err := s.client.SearchLibraries(ctx, input.LibraryName)
			if err != nil {
				return nil, nil, s.toolError(ctx, "failed to retrieve library documentation data from Context7", err)
			}

			out := &resolveLibraryIDOutput{
				Libraries: make([]resolveLibraryIDResult, 0, len(results)),
			}

			if len(results) == 0 {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{
						Text: "No documentation libraries available matching that criteria.",
					}},
				}, out, nil
			}

			content := make([]mcp.Content, 0, len(results))
			var text string

			for _, result := range results {
				lib := resolveLibraryIDResult{
					ResourceURI:   result.GetResourceURI(),
					Title:         result.Title,
					Description:   result.Description,
					TotalSnippets: result.TotalSnippets,
					TrustScore:    result.TrustScore,
				}
				text, err = s.render("resolve_library_id_resp", lib)
				if err != nil {
					return nil, nil, s.toolError(ctx, "failed to render library search results", err)
				}
				content = append(content, &mcp.TextContent{Text: text})
				out.Libraries = append(out.Libraries, lib)
			}
			return &mcp.CallToolResult{Content: content}, out, nil
		},
	}
}
