package tree

type unBranch struct {
	p [fanout]unNode
}

var _ unNode = &unBranch{}

func newUnBranch() *unBranch {
	return &unBranch{}
}

func (b *unBranch) Add(args *CombineArgs, v interface{}, depth int, h hasher, matches *int) unNode {
	i := h.hash()
	n := b.p[i]
	if n == nil {
		n = unEmptyNode{}
	}
	b.p[i] = n.Add(args, v, depth+1, h.next(), matches)
	return b
}

func (b *unBranch) appendTo(dest []interface{}) []interface{} {
	for _, e := range b.p {
		if e != nil {
			if dest = e.appendTo(dest); dest == nil {
				break
			}
		}
	}
	return dest
}

func (b *unBranch) Freeze() node {
	var mask masker
	for i, n := range b.p {
		switch n.(type) {
		case nil, unEmptyNode:
		default:
			mask |= newMasker(i)
		}
	}
	data := make([]node, 0, mask.count())
	for m := mask; m != 0; m = m.next() {
		data = append(data, b.p[m.index()].Freeze())
	}
	return &branch{p: packer{mask: mask, data: data}}
}

func (b *unBranch) Get(args *EqArgs, v interface{}, h hasher) *interface{} {
	if n := b.p[h.hash()]; n != nil {
		return n.Get(args, v, h.next())
	}
	return nil
}

func (b *unBranch) Remove(args *EqArgs, v interface{}, depth int, h hasher, matches *int) unNode {
	i := h.hash()
	if n := b.p[i]; n != nil {
		b.p[i] = b.p[i].Remove(args, v, depth+1, h.next(), matches)
		if _, is := b.p[i].(*unBranch); !is {
			var buf [maxLeafLen]interface{}
			if b := b.appendTo(buf[:]); b != nil {
				l := newUnLeaf()
				l = append(l, b...)
				return &l
			}
		}
	}
	return b
}
