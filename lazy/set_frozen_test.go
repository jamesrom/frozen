package lazy_test

import (
	"testing"

	"github.com/arr-ai/frozen"
	. "github.com/arr-ai/frozen/lazy"
)

func TestSetFrozenEmpty(t *testing.T) {
	t.Parallel()

	f := frozen.Set[any]{}
	s := Frozen(f)

	assertSetOps(t, f, s)

	assertFastIsEmpty(t, s)
	assertFastCountEqual(t, 0, s)
	assertFastCountUpToEqual(t, 0, s, 0)
	assertFastCountUpToEqual(t, 0, s, 1)
	assertFastNotHas(t, s, 3)
	assertFastIsEmpty(t, s.Where(func(_ any) bool { return true }))
	assertFastNotIsEmpty(t, s.With(2))
	assertFastIsEmpty(t, s.Without(2))
	assertFastIsEmpty(t, s.With(2).Without(2))
	assertFastIsEmpty(t, s.Map(func(_ any) any { return 42 }))
	assertFastIsEmpty(t, s.Intersection(Frozen(frozen.NewSet[any](1, 2, 3))))
	assertFastIsEmpty(t, s.Intersection(Frozen(frozen.NewSet[any](1, 2, 3))))
	assertFastIsEmpty(t, s.Powerset())
}

func TestSetFrozenSmall(t *testing.T) {
	t.Parallel()

	f := frozen.NewSet[any](1, 2, 3)
	s := Frozen(f)

	assertSetOps(t, f, s)

	assertFastNotIsEmpty(t, s)
	assertFastCountEqual(t, 3, s)
	assertFastCountUpToEqual(t, 2, s, 2)
	assertFastCountUpToEqual(t, 3, s, 3)
	assertFastCountUpToEqual(t, 3, s, 4)
	assertFastHas(t, s, 3)
	assertFastNotHas(t, s, 4)
	assertFastNotIsEmpty(t, s.With(2))
	assertFastNotIsEmpty(t, s.Without(1).Without(2).Without(4))
	assertFastIsEmpty(t, s.Without(1).Without(2).Without(3))
	assertFastNotIsEmpty(t, s.Map(func(_ any) any { return 42 }))
}
