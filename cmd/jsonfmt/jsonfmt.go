// Command jsonfmt pretty prints the JSON input.
// Arrays with only value types are represented on a single line.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	// Read JSON from stdin or file.
	const indent = "\t"
	var input []byte
	var err error
	if len(os.Args) > 1 {
		input, err = os.ReadFile(os.Args[1])
	} else {
		input, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON into an any.
	var data any
	if err := json.Unmarshal(input, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Format with custom logic.
	output, err := formatJSON(data, indent, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(output)
}

// formatJSON recursively formats the JSON with custom array handling and sorted keys.
func formatJSON(value any, indent string, indentCount int) (string, error) {
	space := strings.Repeat(indent, indentCount)

	switch v := value.(type) {
	case map[string]any:
		var lines []string
		lines = append(lines, "{")

		// Get and sort keys alphabetically.
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// Process sorted keys.
		for _, key := range keys {
			formatted, err := formatJSON(v[key], indent, indentCount+1)
			if err != nil {
				return "", err
			}
			lines = append(lines, fmt.Sprintf("%s%s%q: %s,", space, indent, key, formatted))
		}
		if len(lines) > 1 {
			lines[len(lines)-1] = strings.TrimSuffix(lines[len(lines)-1], ",")
		}
		lines = append(lines, space+"}")
		return strings.Join(lines, "\n"), nil

	case []any:
		// Check if array contains only strings or numbers.
		isSimple := true
		for _, item := range v {
			switch item.(type) {
			case string, float64, int, int64, bool, nil:
				continue
			default:
				isSimple = false
				break
			}
		}

		if isSimple {
			// Compact simple arrays.
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(v); err != nil {
				return "", err
			}
			return strings.TrimSpace(buf.String()), nil
		}

		// Pretty-print complex arrays
		var lines []string
		lines = append(lines, "[")
		for _, item := range v {
			formatted, err := formatJSON(item, indent, indentCount+1)
			if err != nil {
				return "", err
			}
			lines = append(lines, space+indent+formatted+",")
		}
		if len(lines) > 1 {
			lines[len(lines)-1] = strings.TrimSuffix(lines[len(lines)-1], ",")
		}
		lines = append(lines, space+"]")
		return strings.Join(lines, "\n"), nil

	case string:
		return fmt.Sprintf("%q", v), nil
	case float64, int, int64:
		return fmt.Sprintf("%v", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case nil:
		return "null", nil
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}
