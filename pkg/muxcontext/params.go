package muxcontext

import (
	"context"
)

var pathParamsContextKey int
var hostParamsContextKey int

// params is a helper function that retrieves a map of string pairs from the given context using the provided key.
// It returns nil if the context is nil or if no map is associated with the key.
func params(ctx context.Context, key *int) map[string]string {
	if ctx == nil {
		return nil
	}

	h := ctx.Value(key)

	if h != nil {
		return h.(map[string]string)
	}

	return nil
}

// WithPathParams associates the given path parameters with the parent context and returns the resulting context.
func WithPathParams(parent context.Context, params map[string]string) context.Context {
	return context.WithValue(parent, &pathParamsContextKey, params)
}

// PathParams retrieves the associated path parameters from the given context.
func PathParams(ctx context.Context) map[string]string {
	return params(ctx, &pathParamsContextKey)
}

// WithHostParams associates the given host parameters with the parent context and returns the resulting context.
func WithHostParams(parent context.Context, params map[string]string) context.Context {
	return context.WithValue(parent, &hostParamsContextKey, params)
}

// HostParams retrieves the associated host parameters from the given context.
func HostParams(ctx context.Context) map[string]string {
	return params(ctx, &hostParamsContextKey)
}
