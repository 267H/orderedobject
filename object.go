package orderedobject

import (
	"bytes"
	"github.com/json-iterator/go"
)

var json = jsoniter.Config{
	EscapeHTML: false,
}.Froze()

type Object[V any] struct {
	pairs []*pair[V]
	idx   map[string]int
}

type pair[V any] struct {
	k string
	v V
}

func NewObject[V any](capacity int) *Object[V] {
	return &Object[V]{
		pairs: make([]*pair[V], 0, capacity),
		idx:   make(map[string]int, capacity),
	}
}

func (o *Object[V]) Set(key string, value V) {
	if i, ok := o.idx[key]; ok {
		o.pairs[i].v = value
		return
	}
	o.idx[key] = len(o.pairs)
	o.pairs = append(o.pairs, &pair[V]{key, value})
}

func (o *Object[V]) Delete(key string) {
	i, ok := o.idx[key]
	if !ok {
		return
	}
	delete(o.idx, key)
	copy(o.pairs[i:], o.pairs[i+1:])
	o.pairs = o.pairs[:len(o.pairs)-1]
	for j := i; j < len(o.pairs); j++ {
		o.idx[o.pairs[j].k] = j
	}
}

func (o *Object[V]) Has(key string) bool {
	_, ok := o.idx[key]
	return ok
}

func (o *Object[V]) Get(key string) V {
	var zero V
	if i, ok := o.idx[key]; ok {
		return o.pairs[i].v
	}
	return zero
}

func (o *Object[V]) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, p := range o.pairs {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		b.WriteString(p.k)
		b.WriteString(`":`)
		enc := json.BorrowStream(&b)
		enc.WriteVal(p.v)
		enc.Flush()
		json.ReturnStream(enc)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}
