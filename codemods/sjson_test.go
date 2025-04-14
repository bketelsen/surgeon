package codemods

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModifyJSON(t *testing.T) {
	tests := []struct {
		name        string
		action      string
		key         string
		value       string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:     "Set non-existing key",
			action:   "set",
			key:      "foo.bar",
			value:    "baz",
			content:  `{"foo":{}}`,
			expected: `{"foo":{"bar":"baz"}}`,
		},
		{
			name:     "Set existing key",
			action:   "set",
			key:      "foo.bar",
			value:    "baz",
			content:  `{"foo":{"bar":"qux"}}`,
			expected: `{"foo":{"bar":"baz"}}`,
		},
		{
			name:     "Delete key",
			action:   "del",
			key:      "foo.bar",
			content:  `{"foo":{"bar":"baz"}}`,
			expected: `{"foo":{}}`,
		},
		{
			name:        "Unknown action",
			action:      "unknown",
			key:         "foo.bar",
			content:     `{"foo":{"bar":"baz"}}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := modifyJSON(tt.action, tt.key, tt.value, tt.content)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.expected, result)
			}
		})
	}
}
