package muxcontext

import (
	"context"

	"proto.zip/studio/mux/pkg/host"
	"proto.zip/studio/mux/pkg/resource"
)

var hostContextKey int
var resourceContextKey int

func WithHost[RH any, EH any](parent context.Context, h *host.Host[RH, EH]) context.Context {
	if h == nil {
		panic("expected host to not be nil")
	}
	return context.WithValue(parent, &hostContextKey, h)
}

func Host[RH any, EH any](ctx context.Context) *host.Host[RH, EH] {
	if ctx == nil {
		return nil
	}

	h := ctx.Value(&hostContextKey)

	if h != nil {
		return h.(*host.Host[RH, EH])
	}

	return nil
}

func WithResource[RH any](parent context.Context, r *resource.Resource[RH]) context.Context {
	if r == nil {
		panic("expected resource to not be nil")
	}
	return context.WithValue(parent, &resourceContextKey, r)
}

func Resource[RH any](ctx context.Context) *resource.Resource[RH] {
	if ctx == nil {
		return nil
	}

	h := ctx.Value(&resourceContextKey)

	if h != nil {
		return h.(*resource.Resource[RH])
	}

	return nil
}
