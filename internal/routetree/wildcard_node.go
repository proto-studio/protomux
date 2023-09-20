package routetree

import (
	"proto.zip/studio/mux/pkg/tokenizer"
)

// WildcardNode represents a node in the route tree that matches any token.
// It embeds a StandardNode to inherit common node functionalities.
type WildcardNode[H any] struct {
	StandardNode[H]
}

// NewWildcardNode creates and initializes a new WildcardNode.
// It returns the node as an interface of type Node.
func NewWildcardNode[H any]() Node[H] {
	n := &WildcardNode[H]{}
	n.initChildren()
	return n
}

// Match checks if the provided token matches the criteria of the WildcardNode.
// Since it's a wildcard, it always returns true.
func (n *WildcardNode[H]) Match(token tokenizer.Token) bool {
	return true
}

// Equal checks if the provided node is a WildcardNode.
func (n *WildcardNode[H]) Equal(b Node[H]) bool {
	_, ok := b.(*WildcardNode[H])
	return ok
}

// Dynamic indicates if the node represents a dynamic segment in the route tree.
// For WildcardNode, it always returns true as it represents a wildcard segment.
func (n *WildcardNode[H]) Dynamic() bool {
	return true
}
