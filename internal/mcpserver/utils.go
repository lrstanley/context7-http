// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package mcpserver

import (
	"encoding/json"
	"fmt"

	"github.com/lrstanley/context7-http/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
)

// resourceSliceToJSON converts a slice of resources to a slice of mcp.ResourceContents,
// with the MIME type set to application/json and JSON marshalled.
func resourceSliceToJSON[T api.Resource](input []T) (results []mcp.ResourceContents, err error) {
	results = make([]mcp.ResourceContents, len(input))
	var data []byte
	for i, result := range input {
		data, err = json.Marshal(result)
		if err != nil {
			return results, fmt.Errorf("failed to marshal library to json: %w", err)
		}

		results[i] = mcp.TextResourceContents{
			URI:      result.GetResourceURI(),
			MIMEType: "application/json",
			Text:     string(data),
		}
	}
	return results, nil
}
