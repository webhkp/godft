package json

import (
	"encoding/json"
	"fmt"
	"os"
)

// Read a JSON file
func ReadJsonFile[T any](path string) T {
	// read file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}

	// json data
	var result T

	// unmarshall it
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("error:", err)
	}

	return result
}
