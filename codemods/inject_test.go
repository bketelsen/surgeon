package codemods

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInject(t *testing.T) {
	tests := []struct {
		name        string
		where       string
		contents    string
		fileContent []byte
		expected    []byte
		expectError bool
	}{
		{
			name:        "Inject at start",
			where:       "start",
			contents:    "Hello",
			fileContent: []byte("World"),
			expected:    []byte("Hello\nWorld"),
		},
		{
			name:        "Inject at end",
			where:       "end",
			contents:    "Goodbye",
			fileContent: []byte("World"),
			expected:    []byte("World\nGoodbye"),
		},
		{
			name:        "Inject at line",
			where:       "1",
			contents:    "Inserted",
			fileContent: []byte("Line1\nLine2\nLine3"),
			expected:    []byte("Line1\nInserted\nLine2\nLine3\n"),
		},
		{
			name:        "Inject at line",
			where:       "2",
			contents:    "Inserted",
			fileContent: []byte("Line1\nLine2\nLine3"),
			expected:    []byte("Line1\nLine2\nInserted\nLine3\n"),
		},
		{
			name:        "Invalid line number",
			where:       "10",
			contents:    "Out of range",
			fileContent: []byte("Line1\nLine2"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := inject(tt.where, tt.contents, tt.fileContent)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
