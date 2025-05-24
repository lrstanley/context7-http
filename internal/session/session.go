// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package session

import (
	"context"
	"net"
	"net/http"

	"github.com/apex/log"
)

type contextKey string

const sessionIDKey contextKey = "session-id"

// AddIDToContext adds the session ID to the context.
func AddIDToContext(ctx context.Context, r *http.Request) context.Context {
	id := r.RemoteAddr

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		id = host
	}

	return log.NewContext(
		context.WithValue(ctx, sessionIDKey, id),
		log.FromContext(ctx).WithField("session-id", id),
	)
}

// GetIDFromContext returns the session ID from the context, or an empty string if not found.
func GetIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(sessionIDKey).(string)
	if !ok {
		return ""
	}
	return id
}
