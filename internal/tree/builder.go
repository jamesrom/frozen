package tree

import "fmt"

// Builder provides a more efficient way to build nodes incrementally.
type Builder struct {
	t Tree
}

func NewBuilder(capacity int) *Builder {
	return &Builder{}
}

func (b *Builder) Count() int {
	return b.t.count
}

func (b *Builder) Add(args *CombineArgs, v elementT) {
	matches := 0
	root := b.t.MutableRoot()
	h := newHasher(v, 0)
	b.t.root = root.Add(args, v, 0, h, &matches)
	b.t.count += 1 - matches
}

func (b *Builder) Remove(args *EqArgs, v elementT) {
	removed := 0
	root := b.t.MutableRoot()
	h := newHasher(v, 0)
	b.t.root = root.Remove(args, v, 0, h, &removed)
	b.t.count -= removed
}

func (b *Builder) Get(args *EqArgs, el elementT) *elementT {
	return b.t.Get(args, el)
}

func (b *Builder) Finish() Tree {
	t := b.Borrow()
	*b = Builder{}
	return t
}

func (b *Builder) Borrow() Tree {
	return Tree{root: b.t.Root(), count: b.t.count}
}

func (b Builder) String() string {
	return b.Borrow().String()
}

func (b Builder) Format(state fmt.State, verb rune) {
	b.Borrow().Format(state, verb)
}
