// +build !frozen_intf,!frozen_safe_ptr

package tree

import "unsafe"

type node struct {
	b branch
}

func (n noderef) Leaf() *leaf {
	if n.b.isLeaf {
		return (*leaf)(unsafe.Pointer(n))
	}
	return nil
}

func (n noderef) Branch() *branch {
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

func (l *leaf) Node() noderef {
	return (noderef)(unsafe.Pointer(l))
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

func (b *branch) Node() noderef {
	return (noderef)(unsafe.Pointer(b))
}
