// Generated by gen-kv.pl. DO NOT EDIT.
package kvt

import (
	"github.com/arr-ai/frozen/pkg/kv"
)

type unLeaf []kv.KeyValue

var _ unNode = &unLeaf{}

func newUnLeaf() unLeaf {
	return make(unLeaf, 0, maxLeafLen)
}

func (l *unLeaf) Add(args *CombineArgs, v kv.KeyValue, depth int, h hasher, matches *int) unNode {
	for i, e := range *l {
		if args.eq(e, v) {
			*matches++
			(*l)[i] = args.f(e, v)
			return l
		}
	}
	if len(*l) < cap(*l) || depth >= maxTreeDepth {
		*l = append(*l, v)
		return l
	}

	b := newUnBranch()
	for _, e := range *l {
		b.Add(args, e, depth, newHasher(e, depth), matches)
	}
	b.Add(args, v, depth, h, matches)

	return b
}

func (l unLeaf) appendTo(dest []kv.KeyValue) []kv.KeyValue {
	if len(dest)+len(l) > cap(dest) {
		return nil
	}
	return append(dest, l...)
}

func (l unLeaf) Freeze() node {
	ret := make(leaf, 0, len(l))
	ret = append(ret, l...)
	return ret
}

func (l unLeaf) Get(args *EqArgs, v kv.KeyValue, h hasher) *kv.KeyValue {
	for i, e := range l {
		if args.eq(e, v) {
			return &(l)[i]
		}
	}
	return nil
}

func (l *unLeaf) Remove(args *EqArgs, v kv.KeyValue, depth int, h hasher, matches *int) unNode {
	for i, e := range *l {
		if args.eq(e, v) {
			*matches++
			last := len(*l) - 1
			if last == 0 {
				return unEmptyNode{}
			}
			if i < last {
				(*l)[i] = (*l)[last]
			}
			*l = (*l)[:last]
			return l
		}
	}
	return l
}
