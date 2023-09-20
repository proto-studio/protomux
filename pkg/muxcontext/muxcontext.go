// Package muxcontext provides context-related utilities and structures for the mux routing system.
// It aids in managing and retrieving route-specific data from the context during request handling.
package muxcontext

import (
	"context"

	"proto.zip/studio/mux/pkg/host"
	"proto.zip/studio/mux/pkg/resource"
)

var hostContextKey int
var resourceContextKey int

// WithHost associates the given host with the parent context and returns the resulting context.
// It panics if the provided host is nil.
func WithHost[RH any, EH any](parent context.Context, h *host.Host[RH, EH]) context.Context {
	if h == nil {
		panic("expected host to not be nil")
	}
	return context.WithValue(parent, &hostContextKey, h)
}

// Host retrieves the associated host from the given context.
// It returns nil if the context is nil or if no host is associated with it.
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

// WithResource associates the given resource with the parent context and returns the resulting context.
// It panics if the provided resource is nil.
func WithResource[RH any](parent context.Context, r *resource.Resource[RH]) context.Context {
	if r == nil {
		panic("expected resource to not be nil")
	}
	return context.WithValue(parent, &resourceContextKey, r)
}

// Resource retrieves the associated resource from the given context.
// It returns nil if the context is nil or if no resource is associated with it.
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
