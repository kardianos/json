package main

import (
	"bytes"
	"testing"
)

func TestCompareJSON(t *testing.T) {
	tests := []struct {
		name     string
		json1    any
		json2    any
		expected string
		wantDiff bool
	}{
		{
			name: "identical objects",
			json1: map[string]any{
				"a": "foo",
				"b": 42,
			},
			json2: map[string]any{
				"b": 42,
				"a": "foo",
			},
			expected: "",
			wantDiff: false,
		},
		{
			name: "missing key in second",
			json1: map[string]any{
				"name": "Alice",
				"age":  30,
			},
			json2: map[string]any{
				"age": 30,
			},
			expected: `.name: missing in second
`,
			wantDiff: true,
		},
		{
			name: "extra key in second",
			json1: map[string]any{
				"name": "Alice",
			},
			json2: map[string]any{
				"name": "Alice",
				"age":  30,
			},
			expected: `.age: extra in second
`,
			wantDiff: true,
		},
		{
			name: "value mismatch in object",
			json1: map[string]any{
				"name": "Alice",
				"age":  30,
			},
			json2: map[string]any{
				"name": "Bob",
				"age":  30,
			},
			expected: `.name: value mismatch (Alice vs Bob)
`,
			wantDiff: true,
		},
		{
			name:     "identical arrays",
			json1:    []any{"a", "b", 1},
			json2:    []any{"a", "b", 1},
			expected: "",
			wantDiff: false,
		},
		{
			name:  "array order mismatch",
			json1: []any{"a", "b"},
			json2: []any{"b", "a"},
			expected: `[0]: value mismatch (a vs b)
[1]: value mismatch (b vs a)
`,
			wantDiff: true,
		},
		{
			name:  "array length mismatch",
			json1: []any{"a", "b"},
			json2: []any{"a"},
			expected: `: length mismatch (2 vs 1)
`,
			wantDiff: true,
		},
		{
			name: "nested object difference",
			json1: map[string]any{
				"person": map[string]any{
					"name": "Alice",
					"age":  30,
				},
			},
			json2: map[string]any{
				"person": map[string]any{
					"name": "Bob",
					"age":  31,
				},
			},
			expected: `.person.age: value mismatch (30 vs 31)
.person.name: value mismatch (Alice vs Bob)
`,
			wantDiff: true,
		},
		{
			name: "nested array difference",
			json1: map[string]any{
				"items": []any{
					map[string]any{"id": 1},
					map[string]any{"id": 2},
				},
			},
			json2: map[string]any{
				"items": []any{
					map[string]any{"id": 1},
					map[string]any{"id": 3},
				},
			},
			expected: `.items[1].id: value mismatch (2 vs 3)
`,
			wantDiff: true,
		},
		{
			name:  "type mismatch (map vs array)",
			json1: map[string]any{"key": "value"},
			json2: []any{"value"},
			expected: `: type mismatch (map vs []interface {})
`,
			wantDiff: true,
		},
		{
			name:  "type mismatch (array vs string)",
			json1: []any{1, 2},
			json2: "hello",
			expected: `: type mismatch (array vs string)
`,
			wantDiff: true,
		},
		{
			name:  "null vs value",
			json1: nil,
			json2: "not null",
			expected: `: value mismatch (<nil> vs not null)
`,
			wantDiff: true,
		},
		{
			name: "complex nested structure",
			json1: map[string]any{
				"a": []any{
					map[string]any{
						"b": "foo",
						"c": 1,
					},
				},
			},
			json2: map[string]any{
				"a": []any{
					map[string]any{
						"b": "bar",
						"c": 2,
					},
				},
			},
			expected: `.a[0].b: value mismatch (foo vs bar)
.a[0].c: value mismatch (1 vs 2)
`,
			wantDiff: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var diff bool
			compareJSON(&buf, "", tt.json1, tt.json2, &diff)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("compareJSON() output = %q, want %q", got, tt.expected)
			}
			if diff != tt.wantDiff {
				t.Errorf("compareJSON() diff = %v, want %v", diff, tt.wantDiff)
			}
		})
	}
}
