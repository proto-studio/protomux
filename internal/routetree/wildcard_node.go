package routetree

import (
	"proto.zip/studio/mux/pkg/tokenizer"
)

type WildcardNode[H any] struct {
	StandardNode[H]
}

func NewWildcardNode[H any]() Node[H] {
	n := &WildcardNode[H]{}
	n.initChildren()
	return n
}

func (n *WildcardNode[H]) Match(token tokenizer.Token) bool {
	return true
}

func (n *WildcardNode[H]) Equal(b Node[H]) bool {
	_, ok := b.(*WildcardNode[H])
	return ok
}

func (n *WildcardNode[H]) Dynamic() bool {
	return true
}
