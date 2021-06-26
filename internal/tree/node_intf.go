// +build frozen_intf

package tree

import "fmt"

type noderef node

type node interface {
	fmt.Stringer

	Leaf() *leaf
	Branch() *branch

	Add(args *CombineArgs, v elementT, depth int, h hasher, matches *int) noderef
	AppendTo(dest []elementT) []elementT
	Canonical(depth int) noderef
	Combine(args *CombineArgs, n2 noderef, depth int, matches *int) noderef
	Difference(args *EqArgs, n2 noderef, depth int, removed *int) noderef
	Empty() bool
	Equal(args *EqArgs, n2 noderef, depth int) bool
	Get(args *EqArgs, v elementT, h hasher) *elementT
	Intersection(args *EqArgs, n2 noderef, depth int, matches *int) noderef
	Iterator(buf [][]noderef) Iterator
	Reduce(args NodeArgs, depth int, r func(values ...elementT) elementT) elementT
	SubsetOf(args *EqArgs, n2 noderef, depth int) bool
	Transform(args *CombineArgs, depth int, count *int, f func(v elementT) elementT) noderef
	Where(args *WhereArgs, depth int, matches *int) noderef
	With(args *CombineArgs, v elementT, depth int, h hasher, matches *int) noderef
	Without(args *EqArgs, v elementT, depth int, h hasher, matches *int) noderef
	Remove(args *EqArgs, v elementT, depth int, h hasher, matches *int) noderef
}

type leaf struct {
	data []elementT
}

func newLeaf(data ...elementT) *leaf {
	return &leaf{data: data}
}

func (l *leaf) Leaf() *leaf {
	return l
}

func (l *leaf) Branch() *branch {
	return nil
}

func (l *leaf) Node() noderef {
	return l
}

type branch struct {
	p packer
}

func newBranch(p *packer) *branch {
	b := &branch{}
	if p != nil {
		b.p = *p
	}
	return b
}

func (b *branch) Leaf() *leaf {
	return nil
}

func (b *branch) Branch() *branch {
	return b
}

func (b *branch) Node() noderef {
	return b
}
