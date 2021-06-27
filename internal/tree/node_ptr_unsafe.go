// +build !frozen_intf,!frozen_ptr_safe

package tree

import "unsafe"

var (
	// assert(sizeof leaf == sizeof node && alignof leaf == alignof node).
	_ = [-(unsafe.Sizeof(leaf{}) ^ unsafe.Sizeof(node{}))]struct{}{}
	_ = [-(unsafe.Alignof(leaf{}) ^ unsafe.Alignof(node{}))]struct{}{}
)

type node struct {
	b branch
}

func (n *node) Leaf() *leaf {
	if n.b.isLeaf {
		return (*leaf)(unsafe.Pointer(n))
	}
	return nil
}

func (n *node) Branch() *branch {
	if !n.b.isLeaf {
		return &n.b
	}
	return nil
}

type leafBase struct {
	isLeaf bool
	data   []elementT
}

type leaf struct {
	leafBase
	_ [unsafe.Sizeof(branch{}) - unsafe.Sizeof(leafBase{})]byte
}

func newLeaf(data ...elementT) *leaf {
	return &leaf{leafBase: leafBase{isLeaf: true, data: data}}
}

func (l *leaf) Node() *node {
	return (*node)(unsafe.Pointer(l))
}

type branch struct {
	isLeaf bool
	p      packer
}

func newBranch(p *packer) *branch {
	n := &node{}
	if p != nil {
		n.b.p = *p
	}
	return &n.b
}

func (b *branch) Node() *node {
	return (*node)(unsafe.Pointer(b))
}
