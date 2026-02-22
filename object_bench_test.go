package orderedobject

import (
	"fmt"
	"testing"
)

var benchResult []byte

func BenchmarkMarshalJSON_Strings_5(b *testing.B) {
	obj := NewObject[any](5)
	for i := 0; i < 5; i++ {
		obj.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Strings_20(b *testing.B) {
	obj := NewObject[any](20)
	for i := 0; i < 20; i++ {
		obj.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value_longer_string_%d", i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Strings_100(b *testing.B) {
	obj := NewObject[any](100)
	for i := 0; i < 100; i++ {
		obj.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Ints_5(b *testing.B) {
	obj := NewObject[any](5)
	for i := 0; i < 5; i++ {
		obj.Set(fmt.Sprintf("key%d", i), i*1000+i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Ints_100(b *testing.B) {
	obj := NewObject[any](100)
	for i := 0; i < 100; i++ {
		obj.Set(fmt.Sprintf("key%d", i), i*1000+i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Mixed(b *testing.B) {
	obj := NewObject[any](8)
	obj.Set("name", "John Doe")
	obj.Set("age", 30)
	obj.Set("active", true)
	obj.Set("score", 99.5)
	obj.Set("email", "john@example.com")
	obj.Set("count", int64(1000000))
	obj.Set("nothing", nil)
	obj.Set("id", uint64(12345678901234))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Nested(b *testing.B) {
	inner := NewObject[any](3)
	inner.Set("street", "123 Main St")
	inner.Set("city", "New York")
	inner.Set("zip", "10001")

	obj := NewObject[any](4)
	obj.Set("name", "John")
	obj.Set("age", 30)
	obj.Set("address", inner)
	obj.Set("active", true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_EscapedStrings(b *testing.B) {
	obj := NewObject[any](5)
	obj.Set("quote", `He said "hello"`)
	obj.Set("backslash", `path\to\file`)
	obj.Set("newline", "line1\nline2\nline3")
	obj.Set("mixed", "tab\there\nnewline\"quote")
	obj.Set("html", "<script>alert('xss')</script>")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_Empty(b *testing.B) {
	obj := NewObject[any](0)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_TypedString_20(b *testing.B) {
	obj := NewObject[string](20)
	for i := 0; i < 20; i++ {
		obj.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkMarshalJSON_TypedInt_20(b *testing.B) {
	obj := NewObject[int](20)
	for i := 0; i < 20; i++ {
		obj.Set(fmt.Sprintf("key%d", i), i*100)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchResult, _ = obj.MarshalJSON()
	}
}

func BenchmarkSet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		obj := NewObject[any](10)
		for j := 0; j < 10; j++ {
			obj.Set(fmt.Sprintf("key%d", j), j)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	obj := NewObject[any](10)
	for j := 0; j < 10; j++ {
		obj.Set(fmt.Sprintf("key%d", j), j)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = obj.Get("key5")
	}
}
