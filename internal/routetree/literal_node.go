package routetree

import (
	"bytes"

	"proto.zip/studio/mux/pkg/tokenizer"
)

type LiteralNode[H any] struct {
	token []byte
	StandardNode[H]
}

func NewLiteralNode[H any](token tokenizer.Token) Node[H] {
	tokenCopy := make([]byte, len(token))
	copy(tokenCopy, token)

	n := &LiteralNode[H]{
		token: tokenCopy,
	}
	n.initChildren()
	return n
}

func (n *LiteralNode[H]) Match(token tokenizer.Token) bool {
	return bytes.Equal(n.token, token)
}

func (n *LiteralNode[H]) Equal(b Node[H]) bool {
	if litNode, ok := b.(*LiteralNode[H]); ok {
		return bytes.Equal(n.token, litNode.token)
	}
	return false
}

func (n *LiteralNode[H]) Token() []byte {
	return n.token
}
