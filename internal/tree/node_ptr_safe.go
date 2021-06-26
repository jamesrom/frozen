// +build frozen_safe_ptr

package tree

type node struct {
	b      branch
	l      leaf
	isLeaf bool
}

func (n noderef) Leaf() *leaf {
	if n.isLeaf {
		return &n.l
	}
	return nil
}

func (n noderef) Branch() *branch {
	if !n.isLeaf {
		return &n.b
	}
	return nil
}

type leaf struct {
	data []elementT
	n    noderef
}

func newLeaf(data ...elementT) *leaf {
	n := &node{isLeaf: true, l: leaf{data: data}}
	n.l.n = n
	return &n.l
}

func (l *leaf) Node() noderef {
	return l.n
}

type branch struct {
	p packer
	n noderef
}

func newBranch(p *packer) *branch {
	n := &node{}
	if p != nil {
		n.b.p = *p
	}
	n.b.n = n
	return &n.b
}

func (b *branch) Node() noderef {
	return b.n
}
