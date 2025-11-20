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
