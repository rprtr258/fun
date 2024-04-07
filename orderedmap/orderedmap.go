package orderedmap

import (
	"github.com/rprtr258/fun"
	"github.com/rprtr258/fun/iter"
)

// compare uses a less function to determine the ordering of 'a' and 'b'. It returns:
//   - -1 if a < b
//   - 1 if a > b
//   - 0 if a == b
func compare[T any](a, b T, less func(T, T) bool) int {
	switch {
	case less(a, b):
		return -1
	case less(b, a):
		return 1
	default:
		return 0
	}
}

type OrderedMap[K, V any] struct {
	root *node[K, V]
	less func(K, K) bool
}

// New returns empty ordered map
func New[K, V any](less func(K, K) bool) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		less: less,
	}
}

// Put associates 'key' with 'value'
func (t *OrderedMap[K, V]) Put(key K, value V) {
	t.root = t.root.add(key, value, t.less)
}

// Remove removes the value associated with 'key'
func (t *OrderedMap[K, V]) Remove(key K) {
	t.root = t.root.remove(key, t.less)
}

// Get returns the value associated with 'key'
func (t *OrderedMap[K, V]) Get(key K) (V, bool) {
	n := t.root.search(key, t.less)
	if n == nil {
		var v V
		return v, false
	}
	return n.value, true
}

// Iter calls 'yield' on every node in the tree in order
func (t *OrderedMap[K, V]) Iter() iter.Seq[fun.Pair[K, V]] {
	return t.root.each
}

func (t *OrderedMap[K, V]) Kth(k int) (K, bool) {
	if k < 0 || k >= t.root.size() {
		var zero K
		return zero, false
	}

	return t.root.kth(k).key, true
}

func (t *OrderedMap[K, V]) Min() (K, bool) {
	return t.Kth(0)
}

func (t *OrderedMap[K, V]) Max() (K, bool) {
	return t.Kth(t.Size() - 1)
}

// Size returns the number of elements in the tree
func (t *OrderedMap[K, V]) Size() int {
	return t.root.size()
}

type node[K, V any] struct {
	key   K
	value V

	height int
	left   *node[K, V]
	right  *node[K, V]
}

func (n *node[K, V]) add(key K, value V, less func(K, K) bool) *node[K, V] {
	if n == nil {
		return &node[K, V]{
			key:    key,
			value:  value,
			height: 1,
			left:   nil,
			right:  nil,
		}
	}

	switch compare(key, n.key, less) {
	case -1:
		n.left = n.left.add(key, value, less)
	case 1:
		n.right = n.right.add(key, value, less)
	default:
		n.value = value
	}
	return n.rebalanceTree()
}

func (n *node[K, V]) remove(key K, less func(K, K) bool) *node[K, V] {
	if n == nil {
		return nil
	}

	switch compare(key, n.key, less) {
	case -1:
		n.left = n.left.remove(key, less)
	case 1:
		n.right = n.right.remove(key, less)
	default:
		switch {
		case n.left != nil && n.right != nil:
			rightMinNode := n.right.findSmallest()
			n.key = rightMinNode.key
			n.value = rightMinNode.value
			n.right = n.right.remove(rightMinNode.key, less)
		case n.left != nil:
			n = n.left
		case n.right != nil:
			n = n.right
		default:
			n = nil
			return n
		}
	}
	return n.rebalanceTree()
}

func (n *node[K, V]) search(key K, less func(K, K) bool) *node[K, V] {
	if n == nil {
		return nil
	}

	switch compare(key, n.key, less) {
	case -1:
		return n.left.search(key, less)
	case 1:
		return n.right.search(key, less)
	default:
		return n
	}
}

func (n *node[K, V]) each(fn func(kv fun.Pair[K, V]) bool) {
	if n == nil {
		return
	}

	iter.Concat(n.left.each, iter.FromMany(fun.Pair[K, V]{n.key, n.value}), n.right.each)(func(kv fun.Pair[K, V]) bool {
		return fn(kv)
	})
}

func (n *node[K, V]) getHeight() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *node[K, V]) recalculateHeight() {
	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
}

func (n *node[K, V]) rebalanceTree() *node[K, V] {
	if n == nil {
		return n
	}

	n.recalculateHeight()

	switch balanceFactor := n.left.getHeight() - n.right.getHeight(); {
	case balanceFactor <= -2:
		if n.right.left.getHeight() > n.right.right.getHeight() {
			n.right = n.right.rotateRight()
		}
		return n.rotateLeft()
	case balanceFactor >= 2:
		if n.left.right.getHeight() > n.left.left.getHeight() {
			n.left = n.left.rotateLeft()
		}
		return n.rotateRight()
	default:
		return n
	}
}

func (n *node[K, V]) rotateLeft() *node[K, V] {
	newRoot := n.right
	n.right = newRoot.left
	newRoot.left = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (n *node[K, V]) rotateRight() *node[K, V] {
	newRoot := n.left
	n.left = newRoot.right
	newRoot.right = n

	n.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (n *node[K, V]) findSmallest() *node[K, V] {
	if n.left != nil {
		return n.left.findSmallest()
	}

	return n
}

func (n *node[K, V]) size() int {
	if n == nil {
		return 0
	}

	return 1 + n.left.size() + n.right.size()
}

func (n *node[K, V]) kth(k int) *node[K, V] {
	switch ls := n.left.size(); {
	case ls > k:
		return n.left.kth(k)
	case ls == k:
		return n
	default:
		return n.right.kth(k - ls - 1)
	}
}
