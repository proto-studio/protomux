package routetree

import (
	"bytes"

	"proto.zip/studio/mux/pkg/tokenizer"
)

// LiteralNode represents a node in the route tree that matches a specific literal token.
// It embeds a StandardNode to inherit common node functionalities.
type LiteralNode[H any] struct {
	token []byte
	StandardNode[H]
}

// NewLiteralNode creates and initializes a new LiteralNode with the given token.
// It returns the node as an interface of type Node.
func NewLiteralNode[H any](token tokenizer.Token) Node[H] {
	tokenCopy := make([]byte, len(token))
	copy(tokenCopy, token)

	n := &LiteralNode[H]{
		token: tokenCopy,
	}
	n.initChildren()
	return n
}

// Match checks if the provided token matches the token of the LiteralNode.
func (n *LiteralNode[H]) Match(token tokenizer.Token) bool {
	return bytes.Equal(n.token, token)
}

// Equal checks if the provided node is a LiteralNode and if its token matches the token of this LiteralNode.
func (n *LiteralNode[H]) Equal(b Node[H]) bool {
	if litNode, ok := b.(*LiteralNode[H]); ok {
		return bytes.Equal(n.token, litNode.token)
	}
	return false
}

// Token returns the token associated with the LiteralNode.
func (n *LiteralNode[H]) Token() []byte {
	return n.token
}
