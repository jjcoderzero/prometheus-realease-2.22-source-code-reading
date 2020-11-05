package httputil

import (
	"context"
	"net"
	"net/http"

	"github.com/prometheus/prometheus/promql"
)

type pathParam struct{}

// ContextWithPath returns a new context with the given path to be used later
// when logging the query.
func ContextWithPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, pathParam{}, path)
}

// ContextFromRequest returns a new context with identifiers of
// the request to be used later when logging the query.
func ContextFromRequest(ctx context.Context, r *http.Request) context.Context {
	var ip string
	if r.RemoteAddr != "" {
		// r.RemoteAddr has no defined format, so don't return error if we cannot split it into IP:Port.
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	var path string
	if v := ctx.Value(pathParam{}); v != nil {
		path = v.(string)
	}
	return promql.NewOriginContext(ctx, map[string]interface{}{
		"httpRequest": map[string]string{
			"clientIP": ip,
			"method":   r.Method,
			"path":     path,
		},
	})
}
