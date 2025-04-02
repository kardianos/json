// Command jsonfmt pretty prints the JSON input.
// Arrays with only value types are represented on a single line.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Read JSON from stdin or file.
	const indent = "\t"
	var input []byte
	var err error

	var writeToFile string
	w := flag.Bool("w", false, "write back to same file")
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		fn := args[0]
		if *w {
			writeToFile = fn
		}
		input, err = os.ReadFile(fn)
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

	var buf bytes.Buffer
	if err := formatJSON(&buf, data, []byte(indent), 0); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	output := buf.Bytes()
	if len(writeToFile) > 0 {
		err := os.WriteFile(writeToFile, output, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
		return
	}
	os.Stdout.Write(output)
	fmt.Println() // Add trailing newline as in original
}

// formatJSON recursively formats the JSON into the provided buffer.
func formatJSON(buf *bytes.Buffer, value any, indent []byte, indentCount int) error {
	// Pre-allocate space for indentation once per call
	space := make([]byte, indentCount*len(indent))
	for i := 0; i < indentCount; i++ {
		copy(space[i*len(indent):], indent)
	}

	// Fixed buffer for number formatting
	var numBuf [32]byte // Sufficient for most int/float64 values

	switch v := value.(type) {
	case map[string]any:
		buf.WriteByte('{')
		buf.WriteByte('\n')

		// Pre-allocate keys slice with exact capacity
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for i, key := range keys {
			buf.Write(space)
			buf.Write(indent)
			jsonEncodeString(buf, key)
			buf.WriteString(": ")
			if err := formatJSON(buf, v[key], indent, indentCount+1); err != nil {
				return err
			}
			if i < len(keys)-1 {
				buf.WriteByte(',')
			}
			buf.WriteByte('\n')
		}
		buf.Write(space)
		buf.WriteByte('}')

	case []any:
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
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(v); err != nil {
				return err
			}
			// Remove trailing newline from json.Encode
			if buf.Bytes()[buf.Len()-1] == '\n' {
				buf.Truncate(buf.Len() - 1)
			}
		} else {
			buf.WriteByte('[')
			buf.WriteByte('\n')
			for i, item := range v {
				buf.Write(space)
				buf.Write(indent)
				if err := formatJSON(buf, item, indent, indentCount+1); err != nil {
					return err
				}
				if i < len(v)-1 {
					buf.WriteByte(',')
				}
				buf.WriteByte('\n')
			}
			buf.Write(space)
			buf.WriteByte(']')
		}

	case string:
		jsonEncodeString(buf, v)
	case float64:
		// Format float64 into numBuf and write to buf
		n := strconv.AppendFloat(numBuf[:0], v, 'f', -1, 64)
		buf.Write(n)
	case int:
		// Format int into numBuf and write to buf
		n := strconv.AppendInt(numBuf[:0], int64(v), 10)
		buf.Write(n)
	case int64:
		// Format int64 into numBuf and write to buf
		n := strconv.AppendInt(numBuf[:0], v, 10)
		buf.Write(n)
	case bool:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case nil:
		buf.WriteString("null")
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

// jsonEncodeString writes a JSON-quoted string to the buffer.
func jsonEncodeString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch b {
		case '"', '\\':
			buf.WriteByte('\\')
			buf.WriteByte(b)
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		default:
			if b < 32 { // Control characters
				fmt.Fprintf(buf, "\\u00%02x", b) // Fallback for rare cases
			} else {
				buf.WriteByte(b)
			}
		}
	}
	buf.WriteByte('"')
}
