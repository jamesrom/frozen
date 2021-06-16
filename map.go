//nolint:dupl
package frozen

import (
	"fmt"

	"github.com/arr-ai/hash"
)

// KeyValue represents a key-value pair for insertion into a Map.
type KeyValue struct {
	Key, Value interface{}
}

// KV creates a KeyValue.
func KV(key, val interface{}) KeyValue {
	return KeyValue{Key: key, Value: val}
}

// Hash computes a hash for a KeyValue.
func (kv KeyValue) Hash(seed uintptr) uintptr {
	return hash.Interface(kv.Key, seed)
}

func keyValueEqual(a, b interface{}) bool {
	i := a.(KeyValue)
	j := b.(KeyValue)
	return Equal(i.Key, j.Key) && Equal(i.Value, j.Value)
}

// String returns a string representation of a KeyValue.
func (kv KeyValue) String() string {
	return fmt.Sprintf("%#v:%#v", kv.Key, kv.Value)
}

func KeyEqual(a, b interface{}) bool {
	return Equal(a.(KeyValue).Key, b.(KeyValue).Key)
}

// Map maps keys to values. The zero value is the empty Map.
type Map struct {
	root tree
}

var _ Key = Map{}

func newMap(root tree) Map {
	return Map{root: root}
}

// NewMap creates a new Map with kvs as keys and values.
func NewMap(kvs ...KeyValue) Map {
	var b MapBuilder
	for _, kv := range kvs {
		b.Put(kv.Key, kv.Value)
	}
	return b.Finish()
}

// NewMapFromKeys creates a new Map in which values are computed from keys.
func NewMapFromKeys(keys Set, f func(key interface{}) interface{}) Map {
	var b MapBuilder
	for i := keys.Range(); i.Next(); {
		val := i.Value()
		b.Put(val, f(val))
	}
	return b.Finish()
}

// NewMapFromGoMap takes a map[interface{}]interface{} and returns a frozen Map from it.
func NewMapFromGoMap(m map[interface{}]interface{}) Map {
	mb := NewMapBuilder(len(m))
	for k, v := range m {
		mb.Put(k, v)
	}
	return mb.Finish()
}

// IsEmpty returns true if the Map has no entries.
func (m Map) IsEmpty() bool {
	return m.root.count == 0
}

// Count returns the number of entries in the Map.
func (m Map) Count() int {
	return m.root.count
}

// Any returns an arbitrary entry from the Map.
func (m Map) Any() (key, value interface{}) {
	for i := m.Range(); i.Next(); {
		return i.Entry()
	}
	panic("empty map")
}

// With returns a new Map with key associated with val and all other keys
// retained from m.
func (m Map) With(key, val interface{}) Map {
	kv := KV(key, val)
	return newMap(m.root.With(defaultNPKeyCombineArgs, kv))
}

// Without returns a new Map with all keys retained from m except the elements
// of keys.
func (m Map) Without(keys Set) Map {
	args := newEqArgs(
		m.root.Gauge(),
		func(a, b interface{}) bool {
			return Equal(a.(KeyValue).Key, b)
		},
		keyHash,
		hash.Interface)
	return newMap(m.root.Difference(args, keys.Root))
}

// Without2 shoves keys into a Set and calls m.Without.
func (m Map) Without2(keys ...interface{}) Map {
	var sb SetBuilder
	for _, key := range keys {
		sb.Add(key)
	}
	return m.Without(sb.Finish())
}

// Has returns true iff the key exists in the map.
func (m Map) Has(key interface{}) bool {
	return m.root.Get(defaultNPKeyEqArgs, KV(key, nil)) != nil
}

// Get returns the value associated with key in m and true iff the key is found.
func (m Map) Get(key interface{}) (interface{}, bool) {
	if kv := m.root.Get(defaultNPKeyEqArgs, KV(key, nil)); kv != nil {
		return (*kv).(KeyValue).Value, true
	}
	return nil, false
}

// MustGet returns the value associated with key in m or panics if the key is
// not found.
func (m Map) MustGet(key interface{}) interface{} {
	if val, has := m.Get(key); has {
		return val
	}
	panic(fmt.Sprintf("key not found: %v", key))
}

// GetElse returns the value associated with key in m or deflt if the key is not
// found.
func (m Map) GetElse(key, deflt interface{}) interface{} {
	if val, has := m.Get(key); has {
		return val
	}
	return deflt
}

// GetElseFunc returns the value associated with key in m or the result of
// calling deflt if the key is not found.
func (m Map) GetElseFunc(key interface{}, deflt func() interface{}) interface{} {
	if val, has := m.Get(key); has {
		return val
	}
	return deflt()
}

// Keys returns a Set with all the keys in the Map.
func (m Map) Keys() Set {
	var b SetBuilder
	for i := m.Range(); i.Next(); {
		b.Add(i.Key())
	}
	return b.Finish()
}

// Values returns a Set with all the Values in the Map.
func (m Map) Values() Set {
	var b SetBuilder
	for i := m.Range(); i.Next(); {
		b.Add(i.Value())
	}
	return b.Finish()
}

// Project returns a Map with only keys included from this Map.
func (m Map) Project(keys Set) Map {
	return m.Where(func(key, val interface{}) bool {
		return keys.Has(key)
	})
}

// Where returns a Map with only key-value pairs satisfying pred.
func (m Map) Where(pred func(key, val interface{}) bool) Map {
	var b MapBuilder
	for i := m.Range(); i.Next(); {
		if key, val := i.Entry(); pred(key, val) {
			b.Put(key, val)
		}
	}
	return b.Finish()
}

// Map returns a Map with keys from this Map, but the values replaced by the
// result of calling f.
func (m Map) Map(f func(key, val interface{}) interface{}) Map {
	var b MapBuilder
	for i := m.Range(); i.Next(); {
		key, val := i.Entry()
		b.Put(key, f(key, val))
	}
	return b.Finish()
}

// Reduce returns the result of applying f to each key-value pair on the Map.
// The result of each call is used as the acc argument for the next element.
func (m Map) Reduce(f func(acc, key, val interface{}) interface{}, acc interface{}) interface{} {
	for i := m.Range(); i.Next(); {
		acc = f(acc, i.Key(), i.Value())
	}
	return acc
}

func (m Map) eqArgs() *eqArgs {
	return newEqArgs(
		newParallelDepthGauge(m.Count()),
		KeyEqual,
		keyHash,
		keyHash,
	)
}

// Merge returns a map from the merging between two maps, should there be a key overlap,
// the value that corresponds to key will be replaced by the value resulted from the
// provided resolve function.
func (m Map) Merge(n Map, resolve func(key, a, b interface{}) interface{}) Map {
	extractAndResolve := func(a, b interface{}) interface{} {
		i := a.(KeyValue)
		j := b.(KeyValue)
		return KV(i.Key, resolve(i.Key, i.Value, j.Value))
	}
	args := newCombineArgs(m.eqArgs(), extractAndResolve)
	return newMap(m.root.Combine(args, n.root))
}

// Update returns a Map with key-value pairs from n added or replacing existing
// keys.
func (m Map) Update(n Map) Map {
	f := useRHS
	if m.Count() > n.Count() {
		m, n = n, m
		f = useLHS
	}
	return newMap(m.root.Combine(newCombineArgs(m.eqArgs(), f), n.root))
}

// Hash computes a hash val for s.
func (m Map) Hash(seed uintptr) uintptr {
	h := hash.Uintptr(uintptr(3167960924819262823&uint64(^uintptr(0))), seed)
	for i := m.Range(); i.Next(); {
		h ^= hash.Interface(i.Value(), hash.Interface(i.Key(), seed))
	}
	return h
}

// Equal returns true iff i is a Map with all the same key-value pairs as this
// Map.
func (m Map) Equal(i interface{}) bool {
	if n, ok := i.(Map); ok {
		args := newEqArgs(
			newParallelDepthGauge(m.Count()),
			keyValueEqual,
			hash.Interface,
			hash.Interface,
		)
		return m.root.Equal(args, n.root)
	}
	return false
}

// String returns a string representatio of the Map.
func (m Map) String() string {
	return fmt.Sprintf("%v", m)
}

// Format writes a string representation of the Map into state.
func (m Map) Format(state fmt.State, _ rune) {
	state.Write([]byte("("))
	for i, n := m.Range(), 0; i.Next(); n++ {
		if n > 0 {
			state.Write([]byte(", "))
		}
		fmt.Fprintf(state, "%v: %v", i.Key(), i.Value())
	}
	state.Write([]byte(")"))
}

// Range returns a MapIterator over the Map.
func (m Map) Range() *MapIterator {
	return &MapIterator{i: m.root.Iterator()}
}

// MapIterator provides for iterating over a Map.
type MapIterator struct {
	i  Iterator
	kv KeyValue
}

// Next moves to the next key-value pair or returns false if there are no more.
func (i *MapIterator) Next() bool {
	if i.i.Next() {
		var ok bool
		i.kv, ok = i.i.Value().(KeyValue)
		if !ok {
			panic(fmt.Sprintf("Unexpected type: %T", i.i.Value()))
		}
		return true
	}
	return false
}

// Key returns the key for the current entry.
func (i *MapIterator) Key() interface{} {
	return i.kv.Key
}

// Value returns the value for the current entry.
func (i *MapIterator) Value() interface{} {
	return i.kv.Value
}

// Entry returns the current key-value pair as two return values.
func (i *MapIterator) Entry() (key, value interface{}) {
	return i.kv.Key, i.kv.Value
}
