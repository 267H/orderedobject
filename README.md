Example usages:

```go
package main

import (
	"fmt"
	"log"

	"github.com/267H/orderedobject"
)

func main() {
	testBasicJSON()
	testNestedJSON()
	testComplexJSON()
	testSpecialObjectCases()
	fmt.Println("\nOrderedObject tests passed successfully!")
}

func testBasicJSON() {
	fmt.Println("Testing Basic OrderedObject:")
	obj := orderedobject.NewObject[any](3)
	obj.Set("name", "John")
	obj.Set("age", 30)
	obj.Set("city", "New York")

	jsonData, err := obj.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Printf("Basic JSON Output: %s\n", string(jsonData))
	expectedJSON := `{"name":"John","age":30,"city":"New York"}`
	assertEqual(string(jsonData), expectedJSON, "Basic JSON")
	fmt.Println("✓ Basic JSON order and encoding verified")
}

func testNestedJSON() {
	fmt.Println("\nTesting Nested OrderedObject:")

	address := orderedobject.NewObject[any](3)
	address.Set("street", "123 Main St")
	address.Set("city", "New York")
	address.Set("zipcode", "10001")

	contact := orderedobject.NewObject[any](2)
	contact.Set("email", "john@example.com")
	contact.Set("phone", "+1-555-555-5555")

	person := orderedobject.NewObject[any](4)
	person.Set("name", "John Doe")
	person.Set("age", 30)
	person.Set("address", address)
	person.Set("contact", contact)

	jsonData, err := person.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal nested JSON: %v", err)
	}
	fmt.Printf("Nested JSON Output: %s\n", string(jsonData))

	expectedJSON := `{"name":"John Doe","age":30,"address":{"street":"123 Main St","city":"New York","zipcode":"10001"},"contact":{"email":"john@example.com","phone":"+1-555-555-5555"}}`
	assertEqual(string(jsonData), expectedJSON, "Nested JSON")
	fmt.Println("✓ Nested JSON order and encoding verified")
}

func testComplexJSON() {
	fmt.Println("\nTesting Complex OrderedObject:")

	product := orderedobject.NewObject[any](6)
	product.Set("id", "PROD-123")
	product.Set("name", "Super Widget")
	product.Set("price", 99.99)

	variants := make([]interface{}, 2)

	variant1 := orderedobject.NewObject[any](3)
	variant1.Set("color", "red")
	variant1.Set("size", "large")
	variant1.Set("stock", 50)
	variants[0] = variant1

	variant2 := orderedobject.NewObject[any](3)
	variant2.Set("color", "blue")
	variant2.Set("size", "medium")
	variant2.Set("stock", 30)
	variants[1] = variant2

	product.Set("variants", variants)

	metadata := orderedobject.NewObject[any](3)
	metadata.Set("created_at", "2024-01-01")
	metadata.Set("updated_at", "2024-01-02")
	metadata.Set("tags", []string{"new", "featured", "sale"})
	product.Set("metadata", metadata)

	specs := orderedobject.NewObject[any](3)
	dimensions := orderedobject.NewObject[any](3)
	dimensions.Set("length", 10)
	dimensions.Set("width", 5)
	dimensions.Set("height", 2)
	specs.Set("dimensions", dimensions)
	specs.Set("weight", 1.5)
	specs.Set("material", "aluminum")
	product.Set("specifications", specs)

	jsonData, err := product.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal complex JSON: %v", err)
	}
	fmt.Printf("Complex JSON Output: %s\n", string(jsonData))

	expectedJSON := `{"id":"PROD-123","name":"Super Widget","price":99.99,"variants":[{"color":"red","size":"large","stock":50},{"color":"blue","size":"medium","stock":30}],"metadata":{"created_at":"2024-01-01","updated_at":"2024-01-02","tags":["new","featured","sale"]},"specifications":{"dimensions":{"length":10,"width":5,"height":2},"weight":1.5,"material":"aluminum"}}`
	assertEqual(string(jsonData), expectedJSON, "Complex JSON")
	fmt.Println("✓ Complex JSON order and encoding verified")
}

func testSpecialObjectCases() {
	fmt.Println("\nTesting Special OrderedObject Cases:")
	specialObj := orderedobject.NewObject[any](3)
	specialObj.Set("url", "https://example.com/path/to/resource?param=value")
	specialObj.Set("html", "<div>test & demo</div>")
	specialObj.Set("path", "/usr/local/bin")

	specialJSON, err := specialObj.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal JSON with special chars: %v", err)
	}
	expectedSpecialJSON := `{"url":"https://example.com/path/to/resource?param=value","html":"<div>test & demo</div>","path":"/usr/local/bin"}`
	assertEqual(string(specialJSON), expectedSpecialJSON, "Special Characters JSON")
	fmt.Println("✓ Special object cases handled correctly")
}


func assertEqual(got, expected, testName string) {
	if got != expected {
		log.Fatalf("\n%s test failed.\nExpected: %s\nGot: %s", testName, expected, got)
	}
}
