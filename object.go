package orderedobject

import "github.com/json-iterator/go"

var jsonCfg = jsoniter.Config{
	EscapeHTML: false,
}.Froze()

type jsonMarshaler interface {
	MarshalJSON() ([]byte, error)
}

func writeJSONValue(enc *jsoniter.Stream, v any) error {
	switch val := v.(type) {
	case nil:
		enc.WriteNil()
	case string:
		enc.WriteString(val)
	case bool:
		enc.WriteBool(val)
	case int:
		enc.WriteInt(val)
	case int8:
		enc.WriteInt8(val)
	case int16:
		enc.WriteInt16(val)
	case int32:
		enc.WriteInt32(val)
	case int64:
		enc.WriteInt64(val)
	case uint:
		enc.WriteUint(val)
	case uint8:
		enc.WriteUint8(val)
	case uint16:
		enc.WriteUint16(val)
	case uint32:
		enc.WriteUint32(val)
	case uint64:
		enc.WriteUint64(val)
	case float32:
		enc.WriteFloat32(val)
	case float64:
		enc.WriteFloat64(val)
	case jsonMarshaler:
		data, err := val.MarshalJSON()
		if err != nil {
			return err
		}
		if _, err := enc.Write(data); err != nil {
			return err
		}
	default:
		enc.WriteVal(v)
	}

	if enc.Error != nil {
		err := enc.Error
		enc.Error = nil
		return err
	}

	return nil
}

type Object[V any] struct {
	pairs []pair[V]
	idx   map[string]int
}

type pair[V any] struct {
	k string
	v V
}

func NewObject[V any](capacity int) *Object[V] {
	if capacity < 1 {
		capacity = 1
	}
	idxCap := capacity + (capacity >> 1)
	if idxCap < 1 {
		idxCap = 1
	}
	return &Object[V]{
		pairs: make([]pair[V], 0, capacity),
		idx:   make(map[string]int, idxCap),
	}
}

func (o *Object[V]) Set(key string, value V) {
	if o.idx == nil {
		o.idx = make(map[string]int, 1)
	}
	if i, ok := o.idx[key]; ok {
		o.pairs[i].v = value
		return
	}
	o.idx[key] = len(o.pairs)
	o.pairs = append(o.pairs, pair[V]{k: key, v: value})
}

func (o *Object[V]) Delete(key string) {
	i, ok := o.idx[key]
	if !ok {
		return
	}
	delete(o.idx, key)

	pairs := o.pairs
	last := len(pairs) - 1
	if i == last {
		pairs[last] = pair[V]{}
		o.pairs = pairs[:last]
		return
	}
	copy(pairs[i:], pairs[i+1:])
	for j := i; j < last; j++ {
		o.idx[pairs[j].k] = j
	}
	pairs[last] = pair[V]{}
	o.pairs = pairs[:last]
}

func (o *Object[V]) Has(key string) bool {
	_, ok := o.idx[key]
	return ok
}

func (o *Object[V]) Get(key string) V {
	if i, ok := o.idx[key]; ok {
		return o.pairs[i].v
	}
	var zero V
	return zero
}

func (o *Object[V]) MarshalJSON() ([]byte, error) {
	pairs := o.pairs
	n := len(pairs)
	if n == 0 {
		return []byte{'{', '}'}, nil
	}

	enc := jsonCfg.BorrowStream(nil)

	enc.WriteObjectStart()
	for i := 0; i < n; i++ {
		if i > 0 {
			enc.WriteMore()
		}
		p := &pairs[i]
		enc.WriteObjectField(p.k)
		if err := writeJSONValue(enc, p.v); err != nil {
			jsonCfg.ReturnStream(enc)
			return nil, err
		}
	}
	enc.WriteObjectEnd()

	if enc.Error != nil {
		err := enc.Error
		enc.Error = nil
		jsonCfg.ReturnStream(enc)
		return nil, err
	}

	out := append([]byte(nil), enc.Buffer()...)
	jsonCfg.ReturnStream(enc)
	return out, nil
}
