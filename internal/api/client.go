// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/lrstanley/chix"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

const (
	context7BaseURL = "https://context7.com/api"
	maxLibraryCache = 100
)

type Client struct {
	HTTPClient             *http.Client
	limiter                limiter.Store
	searchLibraryCache     *cache.Cache[string, []*SearchResult]
	searchLibraryDocsCache *cache.Cache[string, string]
	listLibrariesCache     *cache.Cache[string, []*Library]
}

// New creates a new API client, with associated rate limiting and caching.
func New(ctx context.Context, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	c := &Client{
		HTTPClient: httpClient,
		searchLibraryCache: cache.NewContext(
			ctx,
			cache.AsLFU[string, []*SearchResult](lfu.WithCapacity(maxLibraryCache)),
		),
		listLibrariesCache: cache.NewContext(
			ctx,
			cache.AsLFU[string, []*Library](lfu.WithCapacity(maxLibraryCache)),
		),
		searchLibraryDocsCache: cache.NewContext(
			ctx,
			cache.AsLFU[string, string](lfu.WithCapacity(maxLibraryCache)),
		),
	}

	limiter, err := memorystore.New(&memorystore.Config{
		Tokens:   10,
		Interval: 60 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create limiter: %w", err)
	}
	c.limiter = limiter

	return c, nil
}

func (c *Client) checkRateLimit(ctx context.Context, namespace string) (err error) {
	ip := chix.GetContextIP(ctx)
	_, _, reset, allowed, _ := c.limiter.Take(ctx, namespace+"/"+ip.String())
	if !allowed {
		return fmt.Errorf("rate limit exceeded (reset in %s)", time.Until(time.Unix(0, int64(reset))))
	}
	return nil
}

type Resource interface {
	GetResourceURI() string
}

// ValidateResourceURI validates a resource URI, and optionally checks that the provided
// type matches the host portion of the URI.
func ValidateResourceURI(uri string, optionalType string) (*url.URL, error) {
	resource, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource URI: %w", err)
	}
	if resource.Scheme != "context7" {
		return nil, fmt.Errorf("invalid resource URI scheme: %s", resource.Scheme)
	}
	if optionalType != "" {
		if resource.Host != optionalType {
			return nil, fmt.Errorf("invalid resource URI type: %s", resource.Host)
		}
	}
	return resource, nil
}
