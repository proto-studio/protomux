// Package host handles the implementation of host entries in the mux.
//
// Hosts may be a fully qualified domain or contain expressions.
// Each host has it's own route tree.
package host

import (
	"fmt"
	"strings"

	"proto.zip/studio/mux/internal/routetree"
	"proto.zip/studio/mux/internal/tokenizers"
	"proto.zip/studio/mux/pkg/resource"
	"proto.zip/studio/mux/pkg/tokenizer"
)

// Host structs represent a host entry in the routing tree and are used to match
// incoming requests for a specific host.
type Host[RequestHandlerType any, ErrorHandlerType any] struct {
	routes       routetree.Node[resource.Resource[RequestHandlerType]]
	params       []tokenizer.Token
	ErrorHandler ErrorHandlerType // The function that is called when an error occurs. Nil will route the errors to the default handler.
}

// New creates a new Host entry with the specific request and error handler types.
//
// You won't normally call this directly unless you are implementing a non-standard mux.
// Most of the time you will want to use NewHost() on the mux implementation instead.
func New[RH any, EH any]() *Host[RH, EH] {
	return &Host[RH, EH]{
		routes: routetree.NewWildcardNode[resource.Resource[RH]](),
	}
}

// NewWithParams creates a new host with pattern parameters.
func NewWithParams[RH any, EH any](params []tokenizer.Token) *Host[RH, EH] {
	return &Host[RH, EH]{
		params: params,
		routes: routetree.NewWildcardNode[resource.Resource[RH]](),
	}
}

// Resource fetches a resource under the host route tree.
// It won't create a new resources. If the path does not match any resources then this method
// will return nil.
//
// On success, it will also return the tokens (if any) that matched the path expressions.
func (h *Host[RH, EH]) Resource(path []byte) (*resource.Resource[RH], []tokenizer.Token) {
	tok := tokenizers.NewPathTokenizer(path)

	var paramValues []tokenizer.Token

	node := h.routes
	token, _, _ := tok.Next()
	for node != nil && token != nil {
		node = node.Child(token)

		if node != nil && node.Dynamic() {
			if paramValues == nil {
				paramValues = make([]tokenizer.Token, 0, 1)
			}
			paramValues = append(paramValues, token)
		}

		token, _, _ = tok.Next()
	}

	if node == nil {
		return nil, nil
	}

	return node.Value(), paramValues
}

// NewResources fetches a resource under the host or returns a new one if the resource does not
// exist yet.
// This method takes a pattern and will return an error if the expressions cannot be parsed.
//
// On success, it will also return the tokens (if any) that matched the path expressions.
func (h *Host[RH, EH]) NewResource(pathPattern []byte) (*resource.Resource[RH], []tokenizer.Token, error) {
	tok := tokenizers.NewPathPatternTokenizer(pathPattern)

	node := h.routes
	token, tokenType, err := tok.Next()
	if err != nil {
		return nil, nil, err
	}

	var paramNames []tokenizer.Token

	for token != nil {
		parent := node

		if tokenType == tokenizer.TokenTypeLabel {
			if paramNames == nil {
				paramNames = make([]tokenizer.Token, 0, 1)
			}
			paramNames = append(paramNames, token)
		}

		node = node.Child(token)
		if node == nil {
			if tokenType == tokenizer.TokenTypeLabel {
				node = routetree.NewWildcardNode[resource.Resource[RH]]()
			} else {
				node = routetree.NewLiteralNode[resource.Resource[RH]](token)
			}
			parent.AddChild(node)
		}

		token, tokenType, err = tok.Next()
		if err != nil {
			return nil, nil, err
		}
	}

	r := node.Value()
	if r == nil {
		r = resource.New[RH]()
		node.SetValue(r)
	}
	return r, paramNames, nil
}

// Handle registers a new resource with the given method and path, associating it with the provided handler.
// It also sets parameter names if any are present in the path.
func (h *Host[RH, EH]) Handle(method, path string, handler RH) {
	resource, paramNames, err := h.NewResource([]byte(path))

	if err != nil {
		panic(err)
	}

	// User supplied input so we convert to upper case for ease of use.
	methodUpper := strings.ToUpper(method)

	if len(paramNames) > 0 {
		resource.SetParamNames(methodUpper, paramNames)
	}

	resource.HandleMethod(methodUpper, handler)
}

// ParamMap maps the provided parameter values to their respective names and returns the resulting map.
// It panics if there's a mismatch between the number of configured parameter names and provided values.
func (h *Host[RH, EH]) ParamMap(paramValues []tokenizer.Token) map[string]string {
	if len(h.params) == 0 {
		return nil
	}

	if len(h.params) != len(paramValues) {
		panic(fmt.Errorf("mismatched parameter length: configured with %d name(s) got %d value(s)", len(h.params), len(paramValues)))
	}

	paramMap := make(map[string]string, len(h.params))

	for paramIdx, paramName := range h.params {
		paramMap[string(paramName)] = string(paramValues[paramIdx])
	}

	return paramMap
}
