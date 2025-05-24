// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/apex/log"
)

// request is a generic function that makes an HTTP request to the given path, with
// the given method, params, and body. If the type of T is a string, the body will be
// read and returned as a string, otherwise [request] will attempt to parse the body
// as JSON.
func request[T any](
	ctx context.Context,
	client *Client,
	method,
	path string,
	params map[string]string,
	body io.Reader,
) (T, error) {
	var result T

	req, err := http.NewRequestWithContext(ctx, method, context7BaseURL+path, body)
	if err != nil {
		return result, fmt.Errorf("failed to initialize request: %w", err)
	}

	// To match that of the node client.
	req.Header.Set("User-Agent", "node")
	req.Header.Set("X-Context7-Source", "mcp-server")

	if params != nil {
		query := req.URL.Query()
		for k, v := range params {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	logger := log.FromContext(ctx).WithFields(log.Fields{
		"method": req.Method,
		"url":    req.URL.String(),
	})

	logger.Info("sending request")
	start := time.Now()
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close() //nolint:errcheck

	logger = logger.WithFields(log.Fields{
		"status":   resp.Status,
		"duration": time.Since(start).Round(time.Millisecond),
	})

	if resp.StatusCode >= 400 {
		logger.Error("request failed")
		return result, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}
	logger.Info("request completed")

	// json decode and wrap in generics. if type of T is string, return the body as a string.
	if _, ok := any(result).(string); ok {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}
		result = any(string(body)).(T)
	} else if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
