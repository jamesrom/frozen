// +build !frozen_intf

package tree

type noderef = *node

func (n noderef) String() string {
	if l := n.Leaf(); l != nil {
		return l.String()
	}
	return n.b.String()
}

func (n noderef) Add(args *CombineArgs, v elementT, depth int, h hasher, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Add(args, v, depth, h, matches)
	}
	return n.b.Add(args, v, depth, h, matches)
}

func (n noderef) AppendTo(dest []elementT) []elementT {
	if l := n.Leaf(); l != nil {
		return l.AppendTo(dest)
	}
	return n.b.AppendTo(dest)
}

func (n noderef) Canonical(depth int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Canonical(depth)
	}
	return n.b.Canonical(depth)
}

func (n noderef) Combine(args *CombineArgs, n2 noderef, depth int, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Combine(args, n2, depth, matches)
	}
	return n.b.Combine(args, n2, depth, matches)
}

func (n noderef) Difference(args *EqArgs, n2 noderef, depth int, removed *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Difference(args, n2, depth, removed)
	}
	return n.b.Difference(args, n2, depth, removed)
}

func (n noderef) Empty() bool {
	if l := n.Leaf(); l != nil {
		return l.Empty()
	}
	return n.b.Empty()
}

func (n noderef) Equal(args *EqArgs, n2 noderef, depth int) bool {
	if l := n.Leaf(); l != nil {
		return l.Equal(args, n2, depth)
	}
	return n.b.Equal(args, n2, depth)
}

func (n noderef) Get(args *EqArgs, v elementT, h hasher) *elementT {
	if l := n.Leaf(); l != nil {
		return l.Get(args, v, h)
	}
	return n.b.Get(args, v, h)
}

func (n noderef) Intersection(args *EqArgs, n2 noderef, depth int, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Intersection(args, n2, depth, matches)
	}
	return n.b.Intersection(args, n2, depth, matches)
}

func (n noderef) Iterator(buf [][]noderef) Iterator {
	if l := n.Leaf(); l != nil {
		return l.Iterator(buf)
	}
	return n.b.Iterator(buf)
}

func (n noderef) Reduce(args NodeArgs, depth int, r func(values ...elementT) elementT) elementT {
	if l := n.Leaf(); l != nil {
		return l.Reduce(args, depth, r)
	}
	return n.b.Reduce(args, depth, r)
}

func (n noderef) SubsetOf(args *EqArgs, n2 noderef, depth int) bool {
	if l := n.Leaf(); l != nil {
		return l.SubsetOf(args, n2, depth)
	}
	return n.b.SubsetOf(args, n2, depth)
}

func (n noderef) Map(args *CombineArgs, depth int, count *int, f func(v elementT) elementT) noderef {
	if l := n.Leaf(); l != nil {
		return l.Map(args, depth, count, f)
	}
	return n.b.Map(args, depth, count, f)
}

func (n noderef) Where(args *WhereArgs, depth int, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Where(args, depth, matches)
	}
	return n.b.Where(args, depth, matches)
}

func (n noderef) With(args *CombineArgs, v elementT, depth int, h hasher, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.With(args, v, depth, h, matches)
	}
	return n.b.With(args, v, depth, h, matches)
}

func (n noderef) Without(args *EqArgs, v elementT, depth int, h hasher, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Without(args, v, depth, h, matches)
	}
	return n.b.Without(args, v, depth, h, matches)
}

func (n noderef) Remove(args *EqArgs, v elementT, depth int, h hasher, matches *int) noderef {
	if l := n.Leaf(); l != nil {
		return l.Remove(args, v, depth, h, matches)
	}
	return n.b.Remove(args, v, depth, h, matches)
}
