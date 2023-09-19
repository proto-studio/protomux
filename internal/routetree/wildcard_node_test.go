package routetree_test

import (
	"testing"

	"proto.zip/studio/mux/internal/routetree"
)

func TestNodeWildcardChildren(t *testing.T) {
	n := routetree.NewWildcardNode[any]()
	NodeStandardAllChildTestHelper(t, n)
}

func TestNodeWildcardMatch(t *testing.T) {
	n := routetree.NewWildcardNode[any]()

	if !n.Match([]byte("a")) {
		t.Error("Expected node wildcard to match a short string")
	}

	if !n.Match([]byte("some longer string")) {
		t.Error("Expected node wildcard to match a longer string")
	}

}

func TestNodeWildcardEqual(t *testing.T) {
	n1 := routetree.NewWildcardNode[any]()
	n2 := routetree.NewWildcardNode[any]()

	if !n1.Equal(n2) {
		t.Error("Expected wildcard node to equal self")
	}

	if !n1.Equal(n2) {
		t.Error("Expected two different wildcard nodes to be equal")
	}
}
