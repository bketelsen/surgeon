package codemods

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceFunction(t *testing.T) {
	originalContent := []byte(`
function foo() {
	echo "Hello, World!"
}

function bar() {
	echo "Goodbye, World!"
}
`)

	replacementContent := []byte(`
function foo() {
	echo "Hello, Universe!"
}
`)

	expectedContent := []byte(`
function foo() {
	echo "Hello, Universe!"
}

function bar() {
	echo "Goodbye, World!"
}
`)

	modifiedContent, err := replaceFunction("foo", replacementContent, originalContent)
	require.NoError(t, err)
	assert.Equal(t, string(expectedContent), string(modifiedContent))
}

func TestReplaceFunction_NotFound(t *testing.T) {
	originalContent := []byte(`
function bar() {
    echo "Goodbye, World!"
}
`)

	replacementContent := []byte(`
function foo() {
    echo "Hello, Universe!"
}
`)

	_, err := replaceFunction("foo", replacementContent, originalContent)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "function \"foo\" not found")
}
