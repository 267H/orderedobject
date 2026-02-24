package orderedobject

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"unsafe"
)

var bufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, 256)
		return &buf
	},
}

const hexDigits = "0123456789abcdef"

func appendJSONString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue
		}
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch c {
		case '"':
			dst = append(dst, '\\', '"')
		case '\\':
			dst = append(dst, '\\', '\\')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hexDigits[c>>4], hexDigits[c&0xf])
		}
		start = i + 1
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	dst = append(dst, '"')
	return dst
}

func appendFloat64(dst []byte, val float64) ([]byte, error) {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return dst, fmt.Errorf("json: unsupported value: %v", val)
	}
	abs := math.Abs(val)
	ff := byte('f')
	if abs != 0 && (abs < 1e-6 || abs >= 1e21) {
		ff = 'e'
	}
	return strconv.AppendFloat(dst, val, ff, -1, 64), nil
}

func appendFloat32(dst []byte, val float32) ([]byte, error) {
	f := float64(val)
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return dst, fmt.Errorf("json: unsupported value: %v", val)
	}
	abs := math.Abs(f)
	ff := byte('f')
	if abs != 0 && (abs < 1e-6 || abs >= 1e21) {
		ff = 'e'
	}
	return strconv.AppendFloat(dst, f, ff, -1, 32), nil
}

func appendJSONValue(dst []byte, v any) ([]byte, error) {
	switch val := v.(type) {
	case nil:
		return append(dst, "null"...), nil
	case string:
		return appendJSONString(dst, val), nil
	case bool:
		if val {
			return append(dst, "true"...), nil
		}
		return append(dst, "false"...), nil
	case int:
		return strconv.AppendInt(dst, int64(val), 10), nil
	case int8:
		return strconv.AppendInt(dst, int64(val), 10), nil
	case int16:
		return strconv.AppendInt(dst, int64(val), 10), nil
	case int32:
		return strconv.AppendInt(dst, int64(val), 10), nil
	case int64:
		return strconv.AppendInt(dst, val, 10), nil
	case uint:
		return strconv.AppendUint(dst, uint64(val), 10), nil
	case uint8:
		return strconv.AppendUint(dst, uint64(val), 10), nil
	case uint16:
		return strconv.AppendUint(dst, uint64(val), 10), nil
	case uint32:
		return strconv.AppendUint(dst, uint64(val), 10), nil
	case uint64:
		return strconv.AppendUint(dst, val, 10), nil
	case float32:
		return appendFloat32(dst, val)
	case float64:
		return appendFloat64(dst, val)
	case json.Marshaler:
		data, err := val.MarshalJSON()
		if err != nil {
			return dst, err
		}
		return append(dst, data...), nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return dst, err
		}
		return append(dst, data...), nil
	}
}

type Object[V any] struct {
	pairs []pair[V]
	idx   map[string]int
}

type pair[V any] struct {
	k  string
	kq string // precomputed `"key":` — trades one alloc per Set() for ~30% faster marshal
	v  V
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

func quoteKey(key string) string {
	buf := appendJSONString(make([]byte, 0, len(key)+3), key)
	buf = append(buf, ':')
	return unsafe.String(&buf[0], len(buf))
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
	o.pairs = append(o.pairs, pair[V]{k: key, kq: quoteKey(key), v: value})
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

// AppendJSON appends the JSON encoding of the object to dst and returns
// the extended buffer. Callers can pool dst to avoid allocations entirely.
func (o *Object[V]) AppendJSON(dst []byte) ([]byte, error) {
	pairs := o.pairs
	n := len(pairs)
	if n == 0 {
		return append(dst, '{', '}'), nil
	}

	dst = append(dst, '{')

	var err error

	// Type-specialized loops avoid interface boxing allocations for typed objects.
	// When V is a concrete type like string or int, any(p.v) would allocate
	// a heap slot per value. Using unsafe.Pointer reads the value directly.
	var zero V
	switch any(zero).(type) {
	case string:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst = appendJSONString(dst, *(*string)(unsafe.Pointer(&p.v)))
		}
	case bool:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			if *(*bool)(unsafe.Pointer(&p.v)) {
				dst = append(dst, "true"...)
			} else {
				dst = append(dst, "false"...)
			}
		}
	case int:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst = strconv.AppendInt(dst, int64(*(*int)(unsafe.Pointer(&p.v))), 10)
		}
	case int64:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst = strconv.AppendInt(dst, *(*int64)(unsafe.Pointer(&p.v)), 10)
		}
	case uint64:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst = strconv.AppendUint(dst, *(*uint64)(unsafe.Pointer(&p.v)), 10)
		}
	case float32:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst, err = appendFloat32(dst, *(*float32)(unsafe.Pointer(&p.v)))
			if err != nil {
				return nil, err
			}
		}
	case float64:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst, err = appendFloat64(dst, *(*float64)(unsafe.Pointer(&p.v)))
			if err != nil {
				return nil, err
			}
		}
	default:
		for i := range n {
			if i > 0 {
				dst = append(dst, ',')
			}
			p := &pairs[i]
			dst = append(dst, p.kq...)
			dst, err = appendJSONValue(dst, any(p.v))
			if err != nil {
				return nil, err
			}
		}
	}

	dst = append(dst, '}')
	return dst, nil
}

func (o *Object[V]) MarshalJSON() ([]byte, error) {
	pairs := o.pairs
	n := len(pairs)
	if n == 0 {
		return []byte{'{', '}'}, nil
	}

	bp := bufPool.Get().(*[]byte)
	buf := (*bp)[:0]

	// Pre-size buffer to avoid repeated growth during encoding.
	if estimated := n*24 + 2; cap(buf) < estimated {
		buf = make([]byte, 0, estimated)
	}

	buf, err := o.AppendJSON(buf)
	if err != nil {
		*bp = buf
		bufPool.Put(bp)
		return nil, err
	}

	out := make([]byte, len(buf))
	copy(out, buf)

	*bp = buf
	bufPool.Put(bp)

	return out, nil
}
