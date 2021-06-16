package frozen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackedWith(t *testing.T) {
	p := packed{}
	for i := 0; i < maxLeafLen; i++ {
		assertEqualPacked(t, p, p.With(newMasker(i), emptyNode{}), i)
	}
	for i := 0; i < maxLeafLen; i++ {
		q := p.With(newMasker(i), leaf{1}).With(newMasker(i), emptyNode{})
		assertEqualPacked(t, p, q, i)
	}
}

func TestPackedWithMulti(t *testing.T) {
	p := packed{}.
		With(newMasker(1), leaf{1, 2}).
		With(newMasker(3), leaf{10, 20}).
		With(newMasker(3), emptyNode{}).
		With(newMasker(5), leaf{3, 4})
	q := packed{}.
		With(newMasker(1), leaf{1, 2}).
		With(newMasker(3), emptyNode{}).
		With(newMasker(5), leaf{3, 4})
	assertEqualPacked(t, p, q)
}

func assertEqualPacked(t *testing.T, expected, actual packed, msgAndArgs ...interface{}) bool {
	t.Helper()

	return assert.Equal(t, expected.mask, actual.mask, msgAndArgs...) &&
		assert.ElementsMatch(t, expected.data, actual.data, msgAndArgs...)
}
