// Command jsondiff compares the logical contents of two JSON files.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s file1.json file2.json\n", os.Args[0])
		os.Exit(1)
	}

	file1, file2 := os.Args[1], os.Args[2]

	// Read and parse first JSON file
	data1, err := os.ReadFile(file1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file1, err)
		os.Exit(1)
	}
	var json1 any
	if err := json.Unmarshal(data1, &json1); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file1, err)
		os.Exit(1)
	}

	// Read and parse second JSON file
	data2, err := os.ReadFile(file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file2, err)
		os.Exit(1)
	}
	var json2 any
	if err := json.Unmarshal(data2, &json2); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file2, err)
		os.Exit(1)
	}

	// Compare and output differences
	var diff bool
	compareJSON(os.Stdout, "$", json1, json2, &diff)
	if diff {
		os.Exit(11)
	}
}

// compareJSON recursively compares two JSON structures and prints differences
func compareJSON(w io.Writer, path string, a, b any, diff *bool) {
	switch va := a.(type) {
	case map[string]any:
		vb, ok := b.(map[string]any)
		if !ok {
			*diff = true
			fmt.Fprintf(w, "%s: type mismatch (map vs %T)\n", path, b)
			return
		}

		// Convert maps to sorted key-value pairs for comparison
		keysA := getSortedKeys(va)
		keysB := getSortedKeys(vb)

		// Check for missing or extra keys
		for _, key := range keysA {
			if _, exists := vb[key]; !exists {
				*diff = true
				fmt.Fprintf(w, "%s.%s: missing in second\n", path, key)
			}
		}
		for _, key := range keysB {
			if _, exists := va[key]; !exists {
				*diff = true
				fmt.Fprintf(w, "%s.%s: extra in second\n", path, key)
			}
		}

		// Compare common keys
		for _, key := range keysA {
			if valB, exists := vb[key]; exists {
				newPath := path + "." + key
				if path == "" {
					newPath = "." + key
				}
				compareJSON(w, newPath, va[key], valB, diff)
			}
		}

	case []any:
		vb, ok := b.([]any)
		if !ok {
			*diff = true
			fmt.Fprintf(w, "%s: type mismatch (array vs %T)\n", path, b)
			return
		}

		// Compare array lengths
		if len(va) != len(vb) {
			*diff = true
			fmt.Fprintf(w, "%s: length mismatch (%d vs %d)\n", path, len(va), len(vb))
		}

		// Compare elements by index
		for i := 0; i < len(va) && i < len(vb); i++ {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			compareJSON(w, newPath, va[i], vb[i], diff)
		}

	case string, float64, int, int64, bool, nil:
		if !reflect.DeepEqual(a, b) {
			*diff = true
			fmt.Fprintf(w, "%s: value mismatch (%v vs %v)\n", path, a, b)
		}

	default:
		*diff = true
		fmt.Fprintf(w, "%s: unsupported type %T\n", path, a)
	}
}

// getSortedKeys returns a sorted slice of keys from a map
func getSortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
