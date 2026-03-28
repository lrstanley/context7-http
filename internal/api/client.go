// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/x/http/utils/httpclog"
	"github.com/lrstanley/x/http/utils/httpcretry"
	cache "github.com/lrstanley/x/sync/cache"
	"github.com/lrstanley/x/sync/cache/policy/lfu"
	"github.com/lrstanley/x/sync/rate"
)

const (
	context7BaseURL = "https://context7.com/api"
	maxLibraryCache = 100
)

type Client struct {
	HTTPClient             *http.Client
	logger                 *slog.Logger
	limiter                *rate.KeyWindowLimiter
	searchLibraryCache     *cache.Cache[string, []*SearchResult]
	searchLibraryDocsCache *cache.Cache[string, string]
}

// New creates a new API client, with associated rate limiting and caching.
func New(ctx context.Context, logger *slog.Logger, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		level := slog.LevelInfo
		httpClient = httpcretry.NewClient(&httpcretry.Config{
			BaseTransport: httpclog.NewTransport(&httpclog.Config{
				Level:         &level,
				Logger:        logger,
				BaseTransport: http.DefaultTransport,
			}),
			MaxRetries:    3,
			RetryCallback: httpcretry.LoggerCallback(logger, slog.LevelWarn),
		})
	}

	c := &Client{
		HTTPClient: httpClient,
		logger:     logger,
		searchLibraryCache: cache.New(
			ctx,
			cache.WithLFU[string, []*SearchResult](lfu.WithCapacity(maxLibraryCache)),
		),
		searchLibraryDocsCache: cache.New(
			ctx,
			cache.WithLFU[string, string](lfu.WithCapacity(maxLibraryCache)),
		),
	}

	c.limiter = rate.NewKeyWindowLimiter(10, 60*time.Second, rate.NewLocalCounter(60*time.Second))

	return c, nil
}

func (c *Client) checkRateLimit(ctx context.Context, namespace string) (err error) {
	ip := chix.GetContextIP(ctx)
	allowed, err := c.limiter.Allow(namespace + "/" + ip.String())
	if err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}
	if !allowed {
		return errors.New("rate limit exceeded")
	}
	return nil
}

type Resource interface {
	GetResourceURI() string
}

// ValidateResourceURI validates a resource URI, and optionally checks that the provided
// type matches the host portion of the URI.
func ValidateResourceURI(uri, optionalType string) (*url.URL, error) {
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
