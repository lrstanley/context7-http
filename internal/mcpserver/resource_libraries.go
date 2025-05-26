// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lrstanley/context7-http/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) resourceLibrariesAll() (resource mcp.Resource, handler server.ResourceHandlerFunc) {
	resource = mcp.NewResource(
		"context7://libraries",
		"get-libraries-all",
		mcp.WithMIMEType("application/json"),
		mcp.WithResourceDescription("Lists all known and tracked libraries."),
	)

	return resource, func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		libraries, err := s.client.ListLibraries(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list libraries: %w", err)
		}
		return resourceSliceToJSON(libraries)
	}
}

func (s *Server) resourceLibrariesTop(top int) (resource mcp.Resource, handler server.ResourceHandlerFunc) {
	resource = mcp.NewResource(
		"context7://libraries/top/"+strconv.Itoa(top),
		"get-libraries-top-"+strconv.Itoa(top),
		mcp.WithMIMEType("application/json"),
		mcp.WithResourceDescription("Lists top "+strconv.Itoa(top)+" libraries, sorted by trust score (if available), otherwise by stars."),
	)

	return resource, func(ctx context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		libraries, err := s.client.ListTopLibraries(ctx, top)
		if err != nil {
			return nil, fmt.Errorf("failed to list top %d libraries: %w", top, err)
		}
		return resourceSliceToJSON(libraries)
	}
}

func (s *Server) resourceLibrary() (template mcp.ResourceTemplate, handler server.ResourceTemplateHandlerFunc) {
	template = mcp.NewResourceTemplate(
		"context7://libraries/{project}",
		"get-library-info",
		mcp.WithTemplateMIMEType("application/json"),
		mcp.WithTemplateDescription("Retrieves information about a specific library."),
	)

	return template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		library, err := s.client.GetLibrary(ctx, request.Params.URI)
		if err != nil {
			return nil, fmt.Errorf("failed to get library: %w", err)
		}
		return resourceSliceToJSON([]*api.Library{library})
	}
}
