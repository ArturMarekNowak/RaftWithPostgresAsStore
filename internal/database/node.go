package database

import (
	"bytes"
	"main/internal"
)

type Node struct {
	Items            [internal.MaxItems]*Item
	Children         [internal.MaxChildren]*Node
	NumberOfItems    int
	NumberOfChildren int
}

func (n *Node) IsLeaf() bool {
	return n.NumberOfChildren == 0
}

func (n *Node) Search(key []byte) (int, bool) {
	low, mid, high := 0, 0, n.NumberOfItems
	for low < high {
		mid = (low + high) / 2
		cmp := bytes.Compare(key, n.Items[mid].Key)
		switch {
		case cmp > 0:
			low = mid + 1
		case cmp < 0:
			high = mid
		default:
			return mid, true
		}
	}
	return low, false
}

func (n *Node) InsertItemAt(pos int, i *Item) {
	if pos < n.NumberOfItems {
		// Make space for insertion if we are not appending to the very end of the "Items" array
		copy(n.Items[pos+1:n.NumberOfItems+1], n.Items[pos:n.NumberOfItems])
	}
	n.Items[pos] = i
	n.NumberOfItems++
}

func (n *Node) InsertChildAt(pos int, c *Node) {
	if pos < n.NumberOfChildren {
		// Make space for insertion if we are not appending to the very end of the "Children" array.
		copy(n.Children[pos+1:n.NumberOfChildren+1], n.Children[pos:n.NumberOfChildren])
	}
	n.Children[pos] = c
	n.NumberOfChildren++
}

func (n *Node) Split() (*Item, *Node) {
	// Retrieve the middle Item.
	mid := internal.MinItems
	midItem := n.Items[mid]

	// Create a new Node and copy half of the Items from the current Node to the new Node.
	newNode := &Node{}
	copy(newNode.Items[:], n.Items[mid+1:])
	newNode.NumberOfItems = internal.MinItems

	// If necessary, copy half of the child pointers from the current Node to the new Node.
	if !n.IsLeaf() {
		copy(newNode.Children[:], n.Children[mid+1:])
		newNode.NumberOfChildren = internal.MinItems + 1
	}

	// Remove data Items and child pointers from the current Node that were moved to the new Node.
	for i, l := mid, n.NumberOfItems; i < l; i++ {
		n.Items[i] = nil
		n.NumberOfItems--

		if !n.IsLeaf() {
			n.Children[i+1] = nil
			n.NumberOfChildren--
		}
	}

	// Return the middle Item and the newly created Node, so we can link them to the parent.
	return midItem, newNode
}

func (n *Node) Insert(Item *Item) bool {
	pos, found := n.Search(Item.Key)

	// The data Item already exists, so just update its value.
	if found {
		n.Items[pos] = Item
		return false
	}

	// We have reached a leaf Node with sufficient capacity to accommodate insertion, so insert the new data Item.
	if n.IsLeaf() {
		n.InsertItemAt(pos, Item)
		return true
	}

	// If the next Node along the path of traversal is already full, split it.
	if n.Children[pos].NumberOfItems >= internal.MaxItems {
		midItem, newNode := n.Children[pos].Split()
		n.InsertItemAt(pos, midItem)
		n.InsertChildAt(pos+1, newNode)
		// We may need to change our direction after promoting the middle Item to the parent, depending on its key.
		switch cmp := bytes.Compare(Item.Key, n.Items[pos].Key); {
		case cmp < 0:
			// The key we are looking for is still smaller than the key of the middle Item that we took from the child,
			// so we can continue following the same direction.
		case cmp > 0:
			// The middle Item that we took from the child has a key that is smaller than the one we are looking for,
			// so we need to change our direction.
			pos++
		default:
			// The middle Item that we took from the child is the Item we are searching for, so just update its value.
			n.Items[pos] = Item
			return true
		}

	}

	return n.Children[pos].Insert(Item)
}

func (n *Node) RemoveItemAt(pos int) *Item {
	removedItem := n.Items[pos]
	n.Items[pos] = nil
	// Fill the gap, if the position we are removing from is not the very last occupied position in the "Items" array.
	if lastPos := n.NumberOfItems - 1; pos < lastPos {
		copy(n.Items[pos:lastPos], n.Items[pos+1:lastPos+1])
		n.Items[lastPos] = nil
	}
	n.NumberOfItems--

	return removedItem
}

func (n *Node) RemoveChildAt(pos int) *Node {
	removedChild := n.Children[pos]
	n.Children[pos] = nil
	// Fill the gap, if the position we are removing from is not the very last occupied position in the "Children" array.
	if lastPos := n.NumberOfChildren - 1; pos < lastPos {
		copy(n.Children[pos:lastPos], n.Children[pos+1:lastPos+1])
		n.Children[lastPos] = nil
	}
	n.NumberOfChildren--

	return removedChild
}

func (n *Node) FillChildAt(pos int) {
	switch {
	// Borrow the right-most Item from the left sibling if the left
	// sibling exists and has more than the minimum number of Items.
	case pos > 0 && n.Children[pos-1].NumberOfItems > internal.MinItems:
		// Establish our left and right nodes.
		left, right := n.Children[pos-1], n.Children[pos]
		// Take the Item from the parent and place it at the left-most position of the right Node.
		copy(right.Items[1:right.NumberOfItems+1], right.Items[:right.NumberOfItems])
		right.Items[0] = n.Items[pos-1]
		right.NumberOfItems++
		// For non-leaf nodes, make the right-most child of the left Node the new left-most child of the right Node.
		if !right.IsLeaf() {
			right.InsertChildAt(0, left.RemoveChildAt(left.NumberOfChildren-1))
		}
		// Borrow the right-most Item from the left Node to replace the parent Item.
		n.Items[pos-1] = left.RemoveItemAt(left.NumberOfItems - 1)
	// Borrow the left-most Item from the right sibling if the right
	// sibling exists and has more than the minimum number of Items.
	case pos < n.NumberOfChildren-1 && n.Children[pos+1].NumberOfItems > internal.MinItems:
		// Establish our left and right nodes.
		left, right := n.Children[pos], n.Children[pos+1]
		// Take the Item from the parent and place it at the right-most position of the left Node.
		left.Items[left.NumberOfItems] = n.Items[pos]
		left.NumberOfItems++
		// For non-leaf nodes, make the left-most child of the right Node the new right-most child of the left Node.
		if !left.IsLeaf() {
			left.InsertChildAt(left.NumberOfChildren, right.RemoveChildAt(0))
		}
		// Borrow the left-most Item from the right Node to replace the parent Item.
		n.Items[pos] = right.RemoveItemAt(0)
	// There are no suitable nodes to borrow Items from, so perform a merge.
	default:
		// If we are at the right-most child pointer, merge the Node with its left sibling.
		// In all other cases, we prefer to merge the Node with its right sibling for simplicity.
		if pos >= n.NumberOfItems {
			pos = n.NumberOfItems - 1
		}
		// Establish our left and right nodes.
		left, right := n.Children[pos], n.Children[pos+1]
		// Borrow an Item from the parent Node and place it at the right-most available position of the left Node.
		left.Items[left.NumberOfItems] = n.RemoveItemAt(pos)
		left.NumberOfItems++
		// Migrate all Items from the right Node to the left Node.
		copy(left.Items[left.NumberOfItems:], right.Items[:right.NumberOfItems])
		left.NumberOfItems += right.NumberOfItems
		// For non-leaf nodes, migrate all applicable Children from the right Node to the left Node.
		if !left.IsLeaf() {
			copy(left.Children[left.NumberOfChildren:], right.Children[:right.NumberOfChildren])
			left.NumberOfChildren += right.NumberOfChildren
		}
		// Remove the child pointer from the parent to the right Node and discard the right Node.
		n.RemoveChildAt(pos + 1)
		right = nil
	}
}

func (n *Node) delete(key []byte, isSeekingSuccessor bool) *Item {
	pos, found := n.Search(key)

	var next *Node

	// We have found a Node holding an Item matching the supplied key.
	if found {
		// This is a leaf Node, so we can simply remove the Item.
		if n.IsLeaf() {
			return n.RemoveItemAt(pos)
		}
		// This is not a leaf Node, so we have to find the inorder successor.
		next, isSeekingSuccessor = n.Children[pos+1], true
	} else {
		next = n.Children[pos]
	}

	// We have reached the leaf Node containing the inorder successor, so remove the successor from the leaf.
	if n.IsLeaf() && isSeekingSuccessor {
		return n.RemoveItemAt(0)
	}

	// We were unable to find an Item matching the given key. Don't do anything.
	if next == nil {
		return nil
	}

	// Continue traversing the tree to find an Item matching the supplied key.
	deletedItem := next.delete(key, isSeekingSuccessor)

	// We found the inorder successor, and we are now back at the internal Node containing the Item
	// matching the supplied key. Therefore, we replace the Item with its inorder successor, effectively
	// deleting the Item from the tree.
	if found && isSeekingSuccessor {
		n.Items[pos] = deletedItem
	}

	// Check if an underflow occurred after we deleted an Item down the tree.
	if next.NumberOfItems < internal.MinItems {
		// Repair the underflow.
		if found && isSeekingSuccessor {
			n.FillChildAt(pos + 1)
		} else {
			n.FillChildAt(pos)
		}
	}

	// Propagate the deleted Item back to the previous stack frame.
	return deletedItem
}
