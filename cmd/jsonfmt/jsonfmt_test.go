package main

import (
	"bytes"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	const indent = "\t"
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "simple object with sorted keys",
			input: map[string]interface{}{
				"b": "bar",
				"a": "foo",
			},
			expected: "{\n\t\"a\": \"foo\",\n\t\"b\": \"bar\"\n}",
			wantErr:  false,
		},
		{
			name:     "simple array of strings",
			input:    []interface{}{"a", "b", "c"},
			expected: `["a","b","c"]`,
			wantErr:  false,
		},
		{
			name:     "simple array of numbers",
			input:    []interface{}{1, 2, 3},
			expected: `[1,2,3]`,
			wantErr:  false,
		},
		{
			name:     "simple array of mixed strings and numbers",
			input:    []interface{}{"a", 1, "b", 2},
			expected: `["a",1,"b",2]`,
			wantErr:  false,
		},
		{
			name: "complex array with objects",
			input: []interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
			},
			expected: "[\n\t{\n\t\t\"id\": 1\n\t},\n\t{\n\t\t\"id\": 2\n\t}\n]",
			wantErr:  false,
		},
		{
			name: "nested object with simple array",
			input: map[string]interface{}{
				"name":    "Alice",
				"numbers": []interface{}{1, 2, 3},
			},
			expected: "{\n\t\"name\": \"Alice\",\n\t\"numbers\": [1,2,3]\n}",
			wantErr:  false,
		},
		{
			name: "nested object with complex array",
			input: map[string]interface{}{
				"name": "Alice",
				"items": []interface{}{
					map[string]interface{}{"x": 1},
					map[string]interface{}{"y": 2},
				},
			},
			expected: "{\n\t\"items\": [\n\t\t{\n\t\t\t\"x\": 1\n\t\t},\n\t\t{\n\t\t\t\"y\": 2\n\t\t}\n\t],\n\t\"name\": \"Alice\"\n}",
			wantErr:  false,
		},
		{
			name:     "empty object",
			input:    map[string]interface{}{},
			expected: "{\n}",
			wantErr:  false,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: "[]",
			wantErr:  false,
		},
		{
			name:     "null value",
			input:    nil,
			expected: "null",
			wantErr:  false,
		},
		{
			name:     "single string",
			input:    "hello",
			expected: `"hello"`,
			wantErr:  false,
		},
		{
			name:     "single number",
			input:    42,
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "single bool",
			input:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:    "unsupported type (struct)",
			input:   struct{ X int }{1},
			wantErr: true,
		},
		{
			name: "mixed array with null",
			input: []interface{}{
				"a",
				nil,
				4.94839456352,
				1.4532,
			},
			expected: `["a",null,4.94839456352,1.4532]`,
			wantErr:  false,
		},
		{
			name: "deeply nested structure",
			input: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{
						map[string]interface{}{
							"d": []interface{}{3, 2},
							"c": []interface{}{1, 2},
						},
					},
				},
			},
			expected: "{\n\t\"a\": {\n\t\t\"b\": [\n\t\t\t{\n\t\t\t\t\"c\": [1,2],\n\t\t\t\t\"d\": [3,2]\n\t\t\t}\n\t\t]\n\t}\n}",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := formatJSON(buf, tt.input, []byte(indent), 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := buf.String()
			if err == nil && got != tt.expected {
				t.Errorf("formatJSON() = %q, want %q", got, tt.expected)
			}
		})
	}
}
