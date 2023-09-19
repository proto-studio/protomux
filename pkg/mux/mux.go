package mux

import (
	"proto.zip/studio/mux/internal/routetree"
	"proto.zip/studio/mux/internal/tokenizers"
	"proto.zip/studio/mux/pkg/host"
	"proto.zip/studio/mux/pkg/tokenizer"
)

type Mux[RequestHandlerType any, ErrorHandlerType any] struct {
	defaultHost *host.Host[RequestHandlerType, ErrorHandlerType]
	hosts       routetree.Node[host.Host[RequestHandlerType, ErrorHandlerType]]
}

// Modifies the mux by adding default internal values.
// Required when creating a new mux. Called automatically by New() and NewHttp()
func (m *Mux[RH, EH]) WithDefaults() *Mux[RH, EH] {
	m.defaultHost = host.New[RH, EH]()
	m.hosts = routetree.NewWildcardNode[host.Host[RH, EH]]()
	return m
}

// Creates a new mux.
// Requies types to be specified. In most cases you will want to use NewHttp()
func New[RH any, EH any]() *Mux[RH, EH] {
	return new(Mux[RH, EH]).WithDefaults()
}

// Returns the default host from the Mux. This host is used if no other hosts match the route.
func (m *Mux[RH, EH]) DefaultHost() *host.Host[RH, EH] {
	return m.defaultHost
}

// Creates a new host in the tree using a host pattern.
// Returns all the tokens that correspond to a named parameter.
//
// Example pattern: {subdomain}.example.com
func (m *Mux[RH, EH]) NewHost(hostPattern string) (*host.Host[RH, EH], error) {
	tok := tokenizers.NewDomainPatternTokenizer([]byte(hostPattern))

	var paramNames []tokenizer.Token

	node := m.hosts
	token, tokenType, err := tok.Next()
	if err != nil {
		return nil, err
	}

	for token != nil {
		parent := node

		if tokenType == tokenizer.TokenTypeLabel {
			if paramNames == nil {
				paramNames = make([]tokenizer.Token, 0, 1)
			}
			paramNames = append(paramNames, token[1:len(token)-1])
		}

		node = node.Child(token)
		if node == nil {
			if tokenType == tokenizer.TokenTypeLabel {
				node = routetree.NewWildcardNode[host.Host[RH, EH]]()
			} else {
				node = routetree.NewLiteralNode[host.Host[RH, EH]](token)
			}
			parent.AddChild(node)
		}

		token, tokenType, err = tok.Next()
		if err != nil {
			return nil, err
		}
	}

	h := node.Value()
	if h == nil {
		h = host.NewWithParams[RH, EH](paramNames)
		node.SetValue(h)
	}
	return h, nil
}

// Returns a host matching the hostname or the default host if none is found.
// This functions expects a fully qualified hostname and will not match patterns.
func (m *Mux[RH, EH]) Host(hostname string) (*host.Host[RH, EH], []tokenizer.Token) {
	tok := tokenizers.NewDomainTokenizer([]byte(hostname))

	var paramValues []tokenizer.Token

	node := m.hosts
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
		return m.defaultHost, paramValues
	}

	h := node.Value()
	if h == nil {
		return m.defaultHost, paramValues
	}
	return h, paramValues
}

// Registers a event handler for a specific HTTP method and and path.
func (m *Mux[RH, EH]) Handle(method, path string, handler RH) {
	m.defaultHost.Handle(method, path, handler)
}
