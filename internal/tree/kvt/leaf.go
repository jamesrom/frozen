// Generated by gen-kv.pl kvt kv.KeyValue. DO NOT EDIT.
package kvt

import (
	"fmt"
	"strings"

	"github.com/arr-ai/frozen/errors"
)

var theEmptyNode = newLeaf().Node()

func newMutableLeaf(data ...elementT) *leaf {
	return newLeaf(append(make([]elementT, 0, maxLeafLen), data...)...)
}

func (l *leaf) Add(args *CombineArgs, v elementT, depth int, h hasher, matches *int) noderef {
	for i, e := range l.data {
		if args.eq(e, v) {
			*matches++
			l.data[i] = args.f(e, v)
			return l.Node()
		}
	}
	if len(l.data) < cap(l.data) || depth >= maxTreeDepth {
		l.data = append(l.data, v)
		return l.Node()
	}

	b := newBranch(nil)
	for _, e := range l.data {
		b.Add(args, e, depth, newHasher(e, depth), matches)
	}
	b.Add(args, v, depth, h, matches)

	return b.Node()
}

func (l *leaf) Canonical(depth int) noderef {
	if len(l.data) <= maxLeafLen || depth*fanoutBits >= 64 {
		return l.Node()
	}
	var matches int
	return newBranch(nil).Combine(DefaultNPCombineArgs, l.Node(), depth, &matches)
}

func (l *leaf) Combine(args *CombineArgs, n noderef, depth int, matches *int) noderef {
	if l.Empty() {
		return n
	}

	l2 := n.Leaf()
	if l2 == nil {
		return n.Combine(args.flip, l.Node(), depth, matches)
	}

	cloned := false
scanning:
	for i, e := range l2.data {
		for j, f := range l.data {
			if args.eq(f, e) {
				if !cloned {
					l = l.clone(0)
					cloned = true
				}
				l.data[j] = args.f(f, e)
				*matches++
				continue scanning
			}
		}
		if len(l.data) < maxLeafLen {
			l = newLeaf(append(l.data, e)...)
		} else {
			return newBranch(nil).
				Combine(args, l.Node(), depth, matches).
				Combine(args, newLeaf(l2.data[i:]...).Node(), depth, matches)
		}
	}
	if len(l.data) > maxLeafLen {
		panic(errors.WTF)
	}
	return l.Canonical(depth)
}

func (l *leaf) AppendTo(dest []elementT) []elementT {
	if len(dest)+len(l.data) > cap(dest) {
		return nil
	}
	return append(dest, l.data...)
}

func (l *leaf) Difference(args *EqArgs, n noderef, depth int, removed *int) noderef {
	ret := newLeaf()
	for _, e := range l.data {
		if n.Get(args.flip, e, newHasher(e, depth)) == nil {
			ret.data = append(ret.data, e)
		} else {
			*removed++
		}
	}
	return ret.Canonical(depth)
}

func (l *leaf) Empty() bool {
	return len(l.data) == 0
}

func (l *leaf) Equal(args *EqArgs, n noderef, depth int) bool {
	if m := n.Leaf(); m != nil {
		if len(l.data) != len(m.data) {
			return false
		}
		for _, e := range l.data {
			if m.Get(args, e, 0) == nil {
				return false
			}
		}
		return true
	}
	return false
}

func (l *leaf) Get(args *EqArgs, v elementT, h hasher) *elementT {
	for i, e := range l.data {
		if args.eq(e, v) {
			return &l.data[i]
		}
	}
	return nil
}

func (l *leaf) Intersection(args *EqArgs, n noderef, depth int, matches *int) noderef {
	ret := newLeaf()
	for _, e := range l.data {
		if n.Get(args, e, newHasher(e, depth)) != nil {
			*matches++
			ret.data = append(ret.data, e)
		}
	}
	return ret.Canonical(depth)
}

func (l *leaf) Iterator([][]noderef) Iterator {
	return newSliceIterator(l.data)
}

func (l *leaf) Reduce(_ NodeArgs, _ int, r func(values ...elementT) elementT) elementT {
	return r(l.data...)
}

func (l *leaf) Remove(args *EqArgs, v elementT, depth int, h hasher, matches *int) noderef {
	for i, e := range l.data {
		if args.eq(e, v) {
			*matches++
			last := len(l.data) - 1
			if last == 0 {
				return newMutableLeaf().Node()
			}
			if i < last {
				l.data[i] = l.data[last]
			}
			l.data = l.data[:last]
			return l.Node()
		}
	}
	return l.Node()
}

func (l *leaf) SubsetOf(args *EqArgs, n noderef, _ int) bool {
	for _, e := range l.data {
		if n.Get(args, e, 0) == nil {
			return false
		}
	}
	return true
}

func (l *leaf) Transform(args *CombineArgs, _ int, counts *int, f func(e elementT) elementT) noderef {
	var nb Builder
	for _, e := range l.data {
		nb.Add(args, f(e))
	}
	t := nb.Finish()
	*counts += t.count
	return t.root
}

func (l *leaf) Where(args *WhereArgs, depth int, matches *int) noderef {
	ret := newLeaf()
	for _, e := range l.data {
		if args.Pred(e) {
			ret.data = append(ret.data, e)
			*matches++
		}
	}
	return ret.Canonical(depth)
}

func (l *leaf) With(args *CombineArgs, v elementT, depth int, h hasher, matches *int) noderef {
	for i, e := range l.data {
		if args.eq(e, v) {
			*matches++
			ret := l.clone(0)
			ret.data[i] = args.f(ret.data[i], v)
			return ret.Node()
		}
	}
	return newLeaf(append(l.data, v)...).Canonical(depth)
}

func (l *leaf) Without(args *EqArgs, v elementT, depth int, h hasher, matches *int) noderef {
	for i, e := range l.data {
		if args.eq(e, v) {
			*matches++
			ret := newLeaf(l.data[:len(l.data)-1]...).clone(0)
			if i != len(ret.data) {
				ret.data[i] = l.data[len(ret.data)]
			}
			return ret.Canonical(depth)
		}
	}
	return l.Node()
}

func (l *leaf) clone(extra int) *leaf {
	return newLeaf(append(make([]elementT, 0, len(l.data)+extra), l.data...)...)
}

func (l *leaf) String() string {
	var b strings.Builder
	b.WriteByte('(')
	for i, e := range l.data {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprint(&b, e)
	}
	b.WriteByte(')')
	return b.String()
}
