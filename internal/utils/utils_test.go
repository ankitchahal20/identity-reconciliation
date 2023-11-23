package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	// Set up test data
	arr := []string{"one", "two", "three"}

	// Test positive case
	result := Contains(arr, "two")
	assert.True(t, result)

	// Test negative case
	result = Contains(arr, "four")
	assert.False(t, result)
}