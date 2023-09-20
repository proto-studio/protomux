// Package resource provides structures and utilities for representing and managing web resources in the routing system.
// It encapsulates the behavior and data associated with individual routes or endpoints.
package resource

import (
	"errors"
	"fmt"

	"proto.zip/studio/mux/pkg/tokenizer"
)

// Resource represents a web resource with associated request handlers and parameter mappings.
// A resource may be associated with more than one request method and handler.
// RequestHandlerType is a generic type representing the handler for a specific method.
type Resource[RequestHandlerType any] struct {
	methods  map[string]RequestHandlerType
	paramMap map[string][]tokenizer.Token
}

// New creates and initializes a new Resource instance.
func New[H any]() *Resource[H] {
	return &Resource[H]{
		methods:  make(map[string]H),
		paramMap: make(map[string][]tokenizer.Token),
	}
}

// Method retrieves the request handler associated with the given method name.
// It returns the handler and a boolean indicating if the handler exists.
func (rh *Resource[H]) Method(methodName string) (H, bool) {
	handler, existing := rh.methods[string(methodName)]
	return handler, existing
}

// Methods returns a list of all method names that have associated request handlers in the Resource.
func (rh *Resource[H]) Methods() []string {
	keys := make([]string, 0, len(rh.methods))
	for k := range rh.methods {
		keys = append(keys, k)
	}
	return keys
}

// HandleMethod associates a request handler with the given method name.
// It panics if the method name already has an associated handler.
func (rh *Resource[H]) HandleMethod(methodName string, handler H) {
	nameStr := methodName
	_, existing := rh.methods[nameStr]

	if existing {
		panic(errors.New("can only be called once per method"))
	}

	rh.methods[nameStr] = handler
}

// SetParamNames sets the parameter names for a specific method.
// It panics if parameter names for the method have already been set.
func (rh *Resource[H]) SetParamNames(methodName string, paramNames []tokenizer.Token) {
	nameStr := string(methodName)
	_, existing := rh.paramMap[nameStr]

	if existing {
		panic(errors.New("can only be called once per method"))
	}

	rh.paramMap[nameStr] = paramNames
}

// ParamMap maps the provided parameter values to their respective names for a given method.
// It panics if there's a mismatch between the number of configured parameter names and provided values.
func (rh *Resource[H]) ParamMap(methodName string, paramValues []tokenizer.Token) map[string]string {
	paramNames, ok := rh.paramMap[string(methodName)]

	if !ok && len(paramValues) == 0 {
		return nil
	}

	if len(paramNames) != len(paramValues) {
		panic(fmt.Errorf("mismatched parameter length: configured with %d name(s) got %d value(s)", len(paramNames), len(paramValues)))
	}

	paramMap := make(map[string]string, len(paramNames))

	for paramIdx, paramName := range paramNames {
		paramMap[string(paramName)] = string(paramValues[paramIdx])
	}

	return paramMap
}
