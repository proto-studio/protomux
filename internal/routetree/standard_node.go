package routetree

import (
	"errors"

	"proto.zip/studio/mux/pkg/tokenizer"
)

// StandardNode represents a common node in the route tree.
// It can have literal children (exact matches) and other types of children (like wildcards or parameters).
// H is a generic type representing the handler associated with the node.
type StandardNode[H any] struct {
	literalChildren  map[string]Node[H]
	allOtherChildren []Node[H]
	handler          *H
}

// initChildren initializes the children maps for the StandardNode.
func (n *StandardNode[H]) initChildren() {
	n.literalChildren = make(map[string]Node[H])
	n.allOtherChildren = make([]Node[H], 0)
}

// Child retrieves a child node that matches the provided token.
// It first checks for literal matches and then checks other types of children.
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

// AddChild adds a child node to the current node.
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

// Value returns the handler associated with the node.
func (n *StandardNode[H]) Value() *H {
	return n.handler
}

// SetValue sets the handler for the node.
// It panics if an attempt is made to set a different handler for a node that already has one.
func (n *StandardNode[H]) SetValue(handler *H) {
	if n.handler != nil && n.handler != handler {
		panic(errors.New("attempting to set handler multiple times for the same route"))
	}
	n.handler = handler
}

// Dynamic indicates if the node represents a dynamic segment in the route tree.
// For StandardNode, it always returns false as it represents static segments.
func (n *StandardNode[H]) Dynamic() bool {
	return false
}
