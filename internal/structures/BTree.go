package structures

import (
	"errors"
	"main/internal"
)

type BTree struct {
	root *Node

	find   func([]byte) ([]byte, error)
	insert func(key, val []byte)
	delete func([]byte) bool
}

func (t *BTree) Find(key []byte) ([]byte, error) {
	for next := t.root; next != nil; {
		pos, found := next.Search(key)

		if found {
			return next.Items[pos].Value, nil
		}

		next = next.Children[pos]
	}

	return nil, errors.New("key not found")
}

func (t *BTree) splitRoot() {
	newRoot := &Node{}
	midItem, newNode := t.root.Split()
	newRoot.InsertItemAt(0, midItem)
	newRoot.InsertChildAt(0, t.root)
	newRoot.InsertChildAt(1, newNode)
	t.root = newRoot
}

func (t *BTree) Insert(key, val []byte) {
	i := &Item{key, val}

	// The tree is empty, so initialize a new node.
	if t.root == nil {
		t.root = &Node{}
	}

	// The tree root is full, so perform a split on the root.
	if t.root.NumberOfItems >= internal.MaxItems {
		t.splitRoot()
	}

	// Begin insertion.
	t.root.Insert(i)
}

func (t *BTree) Delete(key []byte) bool {
	if t.root == nil {
		return false
	}
	deletedItem := t.root.delete(key, false)

	if t.root.NumberOfItems == 0 {
		if t.root.IsLeaf() {
			t.root = nil
		} else {
			t.root = t.root.Children[0]
		}
	}

	if deletedItem != nil {
		return true
	}
	return false
}
