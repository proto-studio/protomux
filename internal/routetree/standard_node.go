package routetree

import (
	"errors"

	"proto.zip/studio/mux/pkg/tokenizer"
)

type StandardNode[H any] struct {
	literalChildren  map[string]Node[H]
	allOtherChildren []Node[H]
	handler          *H
}

func (n *StandardNode[H]) initChildren() {
	n.literalChildren = make(map[string]Node[H])
	n.allOtherChildren = make([]Node[H], 0)
}

func (n *StandardNode[H]) Child(token tokenizer.Token) Node[H] {
	if child, ok := n.literalChildren[string(token)]; ok {
		return child
	}
	for _, child := range n.allOtherChildren {
		if child.Match(token) {
			return child
		}
	}
	return nil
}

func (n *StandardNode[H]) AddChild(child Node[H]) {
	if literal, ok := child.(*LiteralNode[H]); ok {
		key := string(literal.token)

		if _, duplicate := n.literalChildren[key]; duplicate {
			panic(errors.New("duplicate path"))
		}

		n.literalChildren[key] = child
	} else {
		for _, existingChild := range n.allOtherChildren {
			if child.Equal(existingChild) {
				panic(errors.New("duplicate path"))
			}
		}

		n.allOtherChildren = append(n.allOtherChildren, child)
	}
}

func (n *StandardNode[H]) Value() *H {
	return n.handler
}

func (n *StandardNode[H]) SetValue(handler *H) {
	if n.handler != nil && n.handler != handler {
		panic(errors.New("attempting to set handler multiple times for the same route"))
	}
	n.handler = handler
}

func (n *StandardNode[H]) Dynamic() bool {
	return false
}
