package routetree_test

import (
	"bytes"
	"testing"

	"proto.zip/studio/mux/internal/routetree"
)

const literalTestAString string = "test_a"

var literalTestA []byte = []byte(literalTestAString)
var literalTestB []byte = []byte("test_b")

func TestNodeLiteralChildren(t *testing.T) {
	root := routetree.NewLiteralNode[any](literalTestA)
	NodeStandardAllChildTestHelper(t, root)
}

func TestNodeLiteralMatch(t *testing.T) {
	n := routetree.NewLiteralNode[any](literalTestA)

	if !n.Match(literalTestA) {
		t.Error("Expected node literal to match")
	}

	if n.Match(literalTestB) {
		t.Error("Expected node literal to not match")
	}
}

func TestNodeLiteralDoesNotMutateInput(t *testing.T) {
	routetree.NewLiteralNode[any](literalTestA)

	// The byte slice should not have been mutated
	if !bytes.Equal(literalTestA, []byte(literalTestAString)) {
		t.Errorf("Expected '%s' to match '%s'", literalTestA, literalTestAString)
	}
}

func TestNodeLiteralDoesNotMutateValue(t *testing.T) {
	/*
		v := make([]byte, len(literalTestAString))
		copy(v, literalTestAString)
		n := routetree.NewLiteralNode[any](v).(*routetree.LiteralNode[any])
		v[0] = 'X'

		// The byte slice should not have been mutated
		if !bytes.Equal(n.Token(), []byte(literalTestAString)) {
			t.Errorf("Expected '%s' to match '%s'", n.Token(), literalTestAString)
		}
	*/
}
