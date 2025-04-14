package codemods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSedReplace(t *testing.T) {
	tests := []struct {
		name        string
		old         string
		newthing    string
		fileContent []byte
		expected    []byte
	}{
		{
			name:        "Replace single occurrence",
			old:         "foo",
			newthing:    "bar",
			fileContent: []byte("foo baz"),
			expected:    []byte("bar baz"),
		},
		{
			name:        "Replace multiple occurrences",
			old:         "foo",
			newthing:    "bar",
			fileContent: []byte("foo foo foo"),
			expected:    []byte("bar bar bar"),
		},

		{
			name:        "Replace multiple occurrences with newlines",
			old:         "foo",
			newthing:    "bar",
			fileContent: []byte("foo\nfoo foo\nfoo baz"),
			expected:    []byte("bar\nbar bar\nbar baz"),
		},
		{
			name:        "No occurrences",
			old:         "foo",
			newthing:    "bar",
			fileContent: []byte("baz qux"),
			expected:    []byte("baz qux"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sedReplace(tt.old, tt.newthing, tt.fileContent)
			assert.Equal(t, tt.expected, result)
		})
	}
}
