package routetree_test

import (
	"testing"

	"proto.zip/studio/mux/internal/routetree"
	"proto.zip/studio/mux/pkg/tokenizer"
)

type SpyLiteralNode struct {
	routetree.LiteralNode[any]
}

func NodeStandardChildTestHelper(t *testing.T, root routetree.Node[any], matchingToken tokenizer.Token, child routetree.Node[any]) {
	t.Helper()

	if c := root.Child(matchingToken); c != nil {
		t.Error("Expected Child to return nil on empty route")
		return
	}

	root.AddChild(child)

	if c := root.Child(matchingToken); c != child {
		t.Error("Expected Child to return added value")
		return
	}

	// Expect second add to panic
	(func() {
		defer func() { _ = recover() }()
		root.AddChild(child)
		t.Errorf("Expected panic")
	})()
}

func NodeStandarLiteralChildTestHelper(t *testing.T, root routetree.Node[any]) *routetree.LiteralNode[any] {
	t.Helper()

	value := []byte("testLit")
	child := routetree.NewLiteralNode[any](value)

	NodeStandardChildTestHelper(t, root, value, child)
	return child.(*routetree.LiteralNode[any])
}

func NodeStandardWildcardChildTestHelper(t *testing.T, root routetree.Node[any]) {
	t.Helper()

	value := []byte("testWildcard") // Should not match the stirng used for any other test
	child := routetree.NewWildcardNode[any]()

	NodeStandardChildTestHelper(t, root, value, child)
}

func NodeStandardAllChildTestHelper(t *testing.T, root routetree.Node[any]) {

	// Add a literal node first
	literalChild := NodeStandarLiteralChildTestHelper(t, root)

	// Add a wildcard node
	NodeStandardWildcardChildTestHelper(t, root)

	// Check that the literal node is still returned when you try to fetch it
	if c := root.Child(literalChild.Token()); c != literalChild {
		t.Error("Expected Child to return added value")
		return
	}
}
