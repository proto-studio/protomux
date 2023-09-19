package routetree

import "proto.zip/studio/mux/pkg/tokenizer"

type Node[V any] interface {
	Match(token tokenizer.Token) bool
	Child(token tokenizer.Token) Node[V]
	AddChild(node Node[V])
	Value() *V
	SetValue(handler *V)
	Equal(node Node[V]) bool
	Dynamic() bool
}
