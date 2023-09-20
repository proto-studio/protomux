package routetree

import "proto.zip/studio/mux/pkg/tokenizer"

// Node represents an interface for nodes in a route tree.
// It provides methods for matching tokens, managing child nodes, and handling associated values.
// V is a generic type representing the value or handler associated with the node.
type Node[V any] interface {
	Match(token tokenizer.Token) bool    // Match checks if the provided token matches the criteria of the node.
	Child(token tokenizer.Token) Node[V] // Child retrieves a child node that matches the provided token.
	AddChild(node Node[V])               // AddChild adds a child node to the current node.
	Value() *V                           // Value returns the value or handler associated with the node.
	SetValue(handler *V)                 // SetValue sets the value or handler associated with the node.
	Equal(node Node[V]) bool             // Equal checks if the provided node is equivalent to the current node.
	Dynamic() bool                       // Dynamic indicates if the node represents a dynamic segment in the route tree, e.g., a wildcard or parameter.

}
