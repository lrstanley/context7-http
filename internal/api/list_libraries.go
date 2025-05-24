// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
)

const (
	minTopLibraries = 250
	maxTopLibraries = 2500
)

type Library struct {
	Settings *LibrarySettings `json:"settings"`
	Version  *LibraryVersion  `json:"version"`
}

func (l *Library) GetResourceURI() string {
	return "context7://libraries/" + strings.TrimLeft(l.Settings.Project, "/")
}

type LibraryVersion struct {
	LastUpdateDate time.Time `json:"lastUpdateDate"`
	TotalSnippets  int       `json:"totalSnippets"`
	TotalTokens    int       `json:"totalTokens"`

	// Fields that we don't currently need.
	// AverageTokens      float64   `json:"averageTokens"`
	// ErrorCount         int       `json:"errorCount"`
	// ParseDurationMilli int       `json:"parseDuration"`
	// ParseTimestamp     time.Time `json:"parseDate"`
	// SHA                string    `json:"sha"`
	// State              string    `json:"state"`
	// TotalPages         int       `json:"totalPages"`
}

type LibrarySettings struct {
	Branch         string   `json:"branch"`
	Description    string   `json:"description"`
	DocsRepoURL    string   `json:"docsRepoUrl"`
	ExcludeFolders []string `json:"excludeFolders"`
	Folders        []string `json:"folders"`
	Project        string   `json:"project"`
	Stars          int      `json:"stars"`
	Title          string   `json:"title"`
	TrustScore     float64  `json:"trustScore"`
}

// ListLibraries returns all known libraries.
func (c *Client) ListLibraries(ctx context.Context) (results []*Library, err error) {
	results, ok := c.listLibrariesCache.Get("all")
	if ok {
		return results, nil
	}

	err = c.checkRateLimit(ctx, "list-libraries")
	if err != nil {
		return results, err
	}

	results, err = request[[]*Library](
		ctx, c,
		http.MethodGet,
		"/libraries",
		nil,
		http.NoBody,
	)
	if err != nil {
		return results, fmt.Errorf("failed to send request: %w", err)
	}

	c.listLibrariesCache.Set("all", results, cache.WithExpiration(time.Minute*30))
	return results, nil
}

// ListTopLibraries returns the top N libraries, sorted by TrustScore (if available),
// otherwise by Stars. Minimum number of results is 50, maximum is 1000.
func (c *Client) ListTopLibraries(ctx context.Context, top int) (results []*Library, err error) {
	top = min(max(top, minTopLibraries), maxTopLibraries)

	key := "all-top-" + strconv.Itoa(top)

	results, ok := c.listLibrariesCache.Get(key)
	if ok {
		return results, nil
	}

	err = c.checkRateLimit(ctx, "list-top-libraries")
	if err != nil {
		return results, err
	}

	libraries, err := c.ListLibraries(ctx)
	if err != nil {
		return results, err
	}

	// Sort by trust score (if available), then stars, then title.
	slices.SortFunc(libraries, func(a, b *Library) int {
		if a.Settings.TrustScore != b.Settings.TrustScore {
			if a.Settings.TrustScore > b.Settings.TrustScore {
				return -1
			}
			return 1
		}

		if a.Settings.Stars != b.Settings.Stars {
			if a.Settings.Stars > b.Settings.Stars {
				return -1
			}
			return 1
		}

		if a.Settings.Title > b.Settings.Title {
			return -1
		}

		return 1
	})

	results = libraries[:top]
	c.listLibrariesCache.Set(key, results, cache.WithExpiration(time.Minute*30))
	return results, nil
}

// GetLibrary returns a library by its resource URI.
func (c *Client) GetLibrary(ctx context.Context, resourceURI string) (library *Library, err error) {
	libraries, err := c.ListLibraries(ctx)
	if err != nil {
		return nil, err
	}
	for _, library := range libraries {
		if library.GetResourceURI() == resourceURI {
			return library, nil
		}
	}
	return nil, fmt.Errorf("library not found")
}
