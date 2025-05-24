// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
)

const DefaultMinimumDocTokens = 10000

type SearchResp struct {
	Results []*SearchResult `json:"results"`
}

type SearchResult struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	LastUpdate    time.Time `json:"lastUpdateDate"`
	TotalTokens   int       `json:"totalTokens"`
	TotalSnippets int       `json:"totalSnippets"`
	Stars         int       `json:"stars"`
	TrustScore    float64   `json:"trustScore,omitempty"`

	// Fields that we don't currently need.
	// Branch     string `json:"branch"`
	// State      string `json:"state"`
	// TotalPages int    `json:"totalPages"`
	// LastUpdate string `json:"lastUpdate"` // Date only.
}

func (s *SearchResult) GetResourceURI() string {
	return "context7://libraries/" + strings.TrimLeft(s.ID, "/")
}

// SearchLibraries searches the Context7 API for libraries matching the given query.
// It returns a list of search results, sorted by relevance.
func (c *Client) SearchLibraries(ctx context.Context, query string) (results []*SearchResult, err error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return results, nil
	}

	results, ok := c.searchLibraryCache.Get(query)
	if ok {
		return results, nil
	}

	err = c.checkRateLimit(ctx, "search-libraries")
	if err != nil {
		return results, err
	}

	resp, err := request[*SearchResp](
		ctx, c,
		http.MethodGet,
		"/v1/search",
		map[string]string{"query": query},
		http.NoBody,
	)
	if err != nil {
		return results, fmt.Errorf("failed to send request: %w", err)
	}

	c.searchLibraryCache.Set(query, resp.Results, cache.WithExpiration(time.Minute*30))
	return resp.Results, nil
}

type SearchLibraryDocsParams struct {
	Topic   string   `json:"topic"`
	Tokens  int      `json:"tokens"`
	Folders []string `json:"folders"`
}

// SearchLibraryDocsText searches the Context7 API for library documentation text matching
// the given resource URI. Result is formatted as LLM-friendly text.
func (c *Client) SearchLibraryDocsText(
	ctx context.Context,
	resourceURI string,
	params *SearchLibraryDocsParams,
) (results string, err error) {
	var resource *url.URL

	resource, err = ValidateResourceURI(resourceURI, "libraries")
	if err != nil {
		return "", err
	}

	if params == nil {
		params = &SearchLibraryDocsParams{}
	}
	if params.Tokens == 0 {
		params.Tokens = DefaultMinimumDocTokens
	}

	query := map[string]string{
		"type":   "txt", // Supports JSON, but unlikely we'd need at this point.
		"tokens": strconv.Itoa(params.Tokens),
	}
	if params.Topic != "" {
		query["topic"] = params.Topic
	}
	if len(params.Folders) > 0 {
		query["folders"] = strings.Join(params.Folders, ",")
	}

	key := resource.String() + ":" + fmt.Sprintf("%v", query)

	var ok bool
	results, ok = c.searchLibraryDocsCache.Get(key)
	if ok {
		return results, nil
	}

	results, err = request[string](
		ctx, c,
		http.MethodGet,
		"/v1/"+resource.Path,
		query,
		http.NoBody,
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	c.searchLibraryDocsCache.Set(key, results, cache.WithExpiration(time.Minute*5))
	return results, nil
}
