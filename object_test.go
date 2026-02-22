package orderedobject

import (
	"encoding/json"
	"testing"
)

func TestSetGetHasDeleteOrder(t *testing.T) {
	obj := NewObject[int](3)
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("a", 3)

	if !obj.Has("a") || !obj.Has("b") {
		t.Fatalf("Has failed after Set")
	}
	if got := obj.Get("a"); got != 3 {
		t.Fatalf("Get returned %d, want 3", got)
	}

	obj.Delete("a")
	if obj.Has("a") {
		t.Fatalf("key should be deleted")
	}
	if got := obj.Get("a"); got != 0 {
		t.Fatalf("deleted key returned %d", got)
	}

	obj.Set("c", 4)
	obj.Set("d", 5)

	jsonBytes, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(jsonBytes) != `{"b":2,"c":4,"d":5}` {
		t.Fatalf("unexpected order/content: %s", jsonBytes)
	}
}

func TestMarshalJSONVariants(t *testing.T) {
	tests := []struct {
		name      string
		obj       *Object[any]
		expected  string
		expectLen int
	}{
		{
			name: "nested object",
			obj: func() *Object[any] {
				address := NewObject[any](3)
				address.Set("street", "123 Main St")
				address.Set("city", "New York")
				address.Set("zipcode", "10001")

				person := NewObject[any](4)
				person.Set("name", "John Doe")
				person.Set("age", 30)
				person.Set("address", address)
				person.Set("active", true)
				return person
			}(),
			expected:  `{"name":"John Doe","age":30,"address":{"street":"123 Main St","city":"New York","zipcode":"10001"},"active":true}`,
			expectLen: 4,
		},
		{
			name: "special characters",
			obj: func() *Object[any] {
				obj := NewObject[any](3)
				obj.Set("url", "https://example.com/path?param=value")
				obj.Set("html", "<div>test & demo</div>")
				obj.Set("path", "/usr/local/bin")
				return obj
			}(),
			expected:  `{"url":"https://example.com/path?param=value","html":"<div>test & demo</div>","path":"/usr/local/bin"}`,
			expectLen: 3,
		},
		{
			name: "empty",
			obj: func() *Object[any] {
				return NewObject[any](0)
			}(),
			expected:  `{}`,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.obj.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(got) != tt.expected {
				t.Fatalf("unexpected json\n got: %s\nwant: %s", got, tt.expected)
			}

			var generic map[string]any
			if err := json.Unmarshal(got, &generic); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v", err)
			}
			if len(generic) != tt.expectLen {
				t.Fatalf("unexpected map size after round trip: %d", len(generic))
			}
		})
	}
}

func TestMarshalAllPrimitiveTypes(t *testing.T) {
	obj := NewObject[any](20)
	obj.Set("nil", nil)
	obj.Set("string", "hello")
	obj.Set("bool_true", true)
	obj.Set("bool_false", false)
	obj.Set("int", int(42))
	obj.Set("int8", int8(8))
	obj.Set("int16", int16(16))
	obj.Set("int32", int32(32))
	obj.Set("int64", int64(64))
	obj.Set("uint", uint(42))
	obj.Set("uint8", uint8(8))
	obj.Set("uint16", uint16(16))
	obj.Set("uint32", uint32(32))
	obj.Set("uint64", uint64(64))
	obj.Set("float32", float32(3.14))
	obj.Set("float64", float64(2.718))

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatalf("round-trip unmarshal failed: %v\njson: %s", err, got)
	}
	if len(m) != 16 {
		t.Fatalf("expected 16 keys, got %d", len(m))
	}
	if m["nil"] != nil {
		t.Fatalf("nil value should unmarshal as nil")
	}
	if m["string"] != "hello" {
		t.Fatalf("string value mismatch: %v", m["string"])
	}
	if m["bool_true"] != true {
		t.Fatalf("bool_true mismatch")
	}
	if m["bool_false"] != false {
		t.Fatalf("bool_false mismatch")
	}
	if m["int"].(float64) != 42 {
		t.Fatalf("int mismatch: %v", m["int"])
	}
}

func TestMarshalStringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"quote", `say "hello"`, `{"k":"say \"hello\""}`},
		{"backslash", `path\to\file`, `{"k":"path\\to\\file"}`},
		{"newline", "line1\nline2", `{"k":"line1\nline2"}`},
		{"tab", "col1\tcol2", `{"k":"col1\tcol2"}`},
		{"carriage return", "a\rb", `{"k":"a\rb"}`},
		{"backspace", "a\bb", `{"k":"a\u0008b"}`},
		{"formfeed", "a\fb", `{"k":"a\u000cb"}`},
		{"null byte", "a\x00b", `{"k":"a\u0000b"}`},
		{"control char 0x1f", "a\x1fb", `{"k":"a\u001fb"}`},
		{"mixed escapes", "a\"b\\c\nd", `{"k":"a\"b\\c\nd"}`},
		{"no escaping needed", "simple ascii", `{"k":"simple ascii"}`},
		{"empty string", "", `{"k":""}`},
		{"unicode passthrough", "hello \u00e9\u00e8", `{"k":"hello éè"}`},
		{"html unescaped", "<b>bold & \"fun\"</b>", `{"k":"<b>bold & \"fun\"</b>"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewObject[any](1)
			obj.Set("k", tt.input)
			got, err := obj.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(got) != tt.expected {
				t.Fatalf("unexpected json\n got: %s\nwant: %s", got, tt.expected)
			}
			var m map[string]any
			if err := json.Unmarshal(got, &m); err != nil {
				t.Fatalf("round-trip failed: %v\njson: %s", err, got)
			}
			if m["k"] != tt.input {
				t.Fatalf("round-trip value mismatch\n got: %q\nwant: %q", m["k"], tt.input)
			}
		})
	}
}

func TestMarshalNegativeNumbers(t *testing.T) {
	obj := NewObject[any](4)
	obj.Set("neg_int", -42)
	obj.Set("neg_int64", int64(-9223372036854775808))
	obj.Set("neg_float", -3.14)
	obj.Set("zero", 0)

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatalf("round-trip unmarshal failed: %v\njson: %s", err, got)
	}
	if m["neg_int"].(float64) != -42 {
		t.Fatalf("neg_int mismatch: %v", m["neg_int"])
	}
	if m["zero"].(float64) != 0 {
		t.Fatalf("zero mismatch: %v", m["zero"])
	}
}

func TestMarshalSingleElement(t *testing.T) {
	obj := NewObject[any](1)
	obj.Set("only", "value")
	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(got) != `{"only":"value"}` {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestMarshalTypedObject(t *testing.T) {
	obj := NewObject[string](3)
	obj.Set("a", "alpha")
	obj.Set("b", "beta")
	obj.Set("c", "gamma")

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(got) != `{"a":"alpha","b":"beta","c":"gamma"}` {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestMarshalTypedIntObject(t *testing.T) {
	obj := NewObject[int](3)
	obj.Set("x", 100)
	obj.Set("y", 200)
	obj.Set("z", 300)

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(got) != `{"x":100,"y":200,"z":300}` {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestMarshalDeepNesting(t *testing.T) {
	inner := NewObject[any](1)
	inner.Set("value", "deep")

	mid := NewObject[any](1)
	mid.Set("inner", inner)

	outer := NewObject[any](1)
	outer.Set("mid", mid)

	got, err := outer.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	expected := `{"mid":{"inner":{"value":"deep"}}}`
	if string(got) != expected {
		t.Fatalf("unexpected\n got: %s\nwant: %s", got, expected)
	}
}

func TestDeleteMiddlePreservesOrder(t *testing.T) {
	obj := NewObject[int](5)
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)
	obj.Set("d", 4)
	obj.Set("e", 5)

	obj.Delete("c")

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(got) != `{"a":1,"b":2,"d":4,"e":5}` {
		t.Fatalf("order not preserved after middle delete: %s", got)
	}
}

func TestDeleteFirstPreservesOrder(t *testing.T) {
	obj := NewObject[int](3)
	obj.Set("a", 1)
	obj.Set("b", 2)
	obj.Set("c", 3)

	obj.Delete("a")

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(got) != `{"b":2,"c":3}` {
		t.Fatalf("order not preserved after first delete: %s", got)
	}
}

func TestDeleteNonexistent(t *testing.T) {
	obj := NewObject[int](2)
	obj.Set("a", 1)
	obj.Delete("nonexistent")

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(got) != `{"a":1}` {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestGetNonexistent(t *testing.T) {
	obj := NewObject[string](1)
	obj.Set("a", "hello")
	if got := obj.Get("nonexistent"); got != "" {
		t.Fatalf("expected zero value, got %q", got)
	}
}

func TestMarshalLargeObject(t *testing.T) {
	obj := NewObject[any](100)
	for i := 0; i < 100; i++ {
		obj.Set("key_"+string(rune('A'+i%26))+string(rune('0'+i/26)), i)
	}

	got, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatalf("round-trip failed: %v", err)
	}
	if len(m) != 100 {
		t.Fatalf("expected 100 keys, got %d", len(m))
	}
}
