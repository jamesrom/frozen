package frozen_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/arr-ai/frozen"
)

func memoizePrepop[T any](prepare func(n int) T) func(n int) T {
	var lk sync.Mutex
	prepop := map[int]T{}
	return func(n int) T {
		lk.Lock()
		defer lk.Unlock()
		if data, has := prepop[n]; has {
			return data
		}
		data := prepare(n)
		prepop[n] = data
		return data
	}
}

// func intSetToMap(s Set) map[int]bool {
// 	result := make(map[int]bool, s.Count())
// 	for i := s.Range(); i.Next(); {
// 		result[i.Value().(int)] = true
// 	}
// 	return result
// }

// func intSetDiff(a, b Set) (l, m, r []int) {
// 	ma := intSetToMap(a)
// 	mb := intSetToMap(b)
// 	for e := range ma {
// 		if mb[e] {
// 			m = append(m, e)
// 		} else {
// 			l = append(l, e)
// 		}
// 	}
// 	for e := range mb {
// 		if !ma[e] {
// 			r = append(r, e)
// 		}
// 	}
// 	sort.Ints(l)
// 	sort.Ints(m)
// 	sort.Ints(r)
// 	return
// }

func assertSetHas[T any](t *testing.T, s Set[T], i T) bool {
	t.Helper()

	return assert.True(t, s.Has(i), "i=%v", i)
}

func assertSetNotHas[T any](t *testing.T, s Set[T], i T) bool {
	t.Helper()

	return assert.False(t, s.Has(i), "i=%v", i)
}

func assertMapEqual[K, V any](
	t *testing.T,
	expected, actual Map[K, V],
	msgAndArgs ...any,
) bool {
	t.Helper()

	format := "\nexpected %v != \nactual   %v"
	args := []any{}
	if len(msgAndArgs) > 0 {
		format = msgAndArgs[0].(string) + format
		args = append(append(args, format), msgAndArgs[1:]...)
	} else {
		args = append(args, format)
	}
	args = append(args, expected, actual)
	return assert.True(t, expected.Equal(actual), args...)
}

func assertMapHas[K, V any](t *testing.T, m Map[K, V], i K, expected V) bool {
	t.Helper()

	v, has := m.Get(i)
	ok1 := assert.Equal(t, has, m.Has(i), "has != Has(): i=%v", i)
	ok2 := assert.True(t, has, "!has: i=%v", i) &&
		assert.Equal(t, expected, v, "expected %v != actual %v: i=%v", expected, v, i)
	return ok1 && ok2
}

func assertMapNotHas[K, V any](t *testing.T, m Map[K, V], i K) bool {
	t.Helper()

	v, has := m.Get(i)
	ok1 := assert.Equal(t, has, m.Has(i))
	ok2 := assert.False(t, has, "i=%v v=%v", i, v)
	return ok1 && ok2
}

func generateSortedIntArray(start, end, step int) []int {
	if step == 0 {
		if start == step {
			return []int{}
		}
		panic("zero step size")
	}
	if (step > 0 && start > end) || (step < 0 && start < end) {
		panic("array will never reach end value")
	}
	n := (start - end) / step
	if n < 0 {
		n *= -1
	}
	result := make([]int, n)
	currentVal := start
	for i := 0; i < n; i++ {
		result[i] = currentVal + step*i
	}
	return result
}
