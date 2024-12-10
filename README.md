Example usages:

```package main

import (
	"fmt"
	"log"

	"github.com/267H/orderedform"
	"github.com/267H/orderedobject"
	"github.com/267H/orderedquery"
)

func main() {
	testBasicJSON()
	testNestedJSON()
	testComplexJSON()
	testOrderedQuery()
	testOrderedForm()
	testSpecialCases()
	fmt.Println("\nAll tests passed successfully!")
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

func testOrderedQuery() {
	fmt.Println("\nTesting OrderedQuery:")

	query := orderedquery.NewQuery(5)
	query.Add("page", "1")
	query.Add("sort", "desc")
	query.Add("filter", "active")
	query.Add("view", "detailed")
	query.Add("limit", "50")

	queryString := query.Encode()
	fmt.Printf("Basic Query string: %s\n", queryString)
	expectedQuery := "page=1&sort=desc&filter=active&view=detailed&limit=50"
	assertEqual(queryString, expectedQuery, "Basic Query")

	complexQuery := orderedquery.NewQuery(4)
	complexQuery.Add("search", "test product")
	complexQuery.Add("category", "electronics/phones")
	complexQuery.Add("price_range", "100-500")
	complexQuery.Add("brand[]", "apple,samsung,google")

	complexQueryString := complexQuery.Encode()
	fmt.Printf("Complex Query string: %s\n", complexQueryString)
	expectedComplexQuery := "search=test+product&category=electronics%2Fphones&price_range=100-500&brand%5B%5D=apple%2Csamsung%2Cgoogle"
	assertEqual(complexQueryString, expectedComplexQuery, "Complex Query")
	fmt.Println("✓ Query strings verified")
}

func testOrderedForm() {
	fmt.Println("\nTesting OrderedForm:")

	form := orderedform.NewForm(6)
	form.Set("username", "john_doe")
	form.Set("email", "john@example.com")
	form.Set("password", "p@ssw0rd!")
	form.Set("confirm_password", "p@ssw0rd!")
	form.Set("preferences[theme]", "dark")
	form.Set("preferences[notifications]", "email,sms")

	formData := form.URLEncode()
	fmt.Printf("Form data: %s\n", formData)
	expectedForm := "username=john_doe&email=john%40example.com&password=p%40ssw0rd%21&confirm_password=p%40ssw0rd%21&preferences%5Btheme%5D=dark&preferences%5Bnotifications%5D=email%2Csms"
	assertEqual(formData, expectedForm, "Complex Form")
	fmt.Println("✓ Form encoding verified")
}

func testSpecialCases() {
	fmt.Println("\nTesting Special Cases:")

	emptyQuery, err := orderedquery.ParseQuery("")
	if err != nil {
		log.Fatalf("Failed to parse empty query: %v", err)
	}
	assertEqual(emptyQuery.Encode(), "", "Empty Query")

	prefixQuery, err := orderedquery.ParseQuery("?key=value&special=a/b/c")
	if err != nil {
		log.Fatalf("Failed to parse query with prefix: %v", err)
	}
	assertEqual(prefixQuery.Encode(), "key=value&special=a%2Fb%2Fc", "Prefixed Query")

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

	fmt.Println("✓ Special cases handled correctly")
}

func assertEqual(got, expected, testName string) {
	if got != expected {
		log.Fatalf("\n%s test failed.\nExpected: %s\nGot: %s", testName, expected, got)
	}
}
