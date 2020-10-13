package utils

import (
	"bytes"
	"fmt"
)

// Please note that I used examples here to assert output
// because its easier to read than a heavily escaped string

func ExamplePrettyJSON_output() {
	data := bytes.NewBufferString("{\"cat\":\"meow\"}").Bytes()
	result := PrettyJSON(data)

	fmt.Println(result)

	// Output:
	// {
	//   "cat": "meow"
	// }
}

func ExamplePrettyJSON_notjson() {
	data := bytes.NewBufferString("<garbage>").Bytes()
	result := PrettyJSON(data)

	fmt.Println(result)

	// Output:
	//
}

func ExamplePrettyJSON_empty() {
	data := bytes.NewBufferString("").Bytes()
	result := PrettyJSON(data)

	fmt.Println(result)

	// Output:
	//
}
