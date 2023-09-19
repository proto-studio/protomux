package muxcontext

import (
	"context"
)

var pathParamsContextKey int
var hostParamsContextKey int

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

func WithPathParams(parent context.Context, params map[string]string) context.Context {
	return context.WithValue(parent, &pathParamsContextKey, params)
}

func PathParams(ctx context.Context) map[string]string {
	return params(ctx, &pathParamsContextKey)
}

func WithHostParams(parent context.Context, params map[string]string) context.Context {
	return context.WithValue(parent, &hostParamsContextKey, params)
}

func HostParams(ctx context.Context) map[string]string {
	return params(ctx, &hostParamsContextKey)
}
