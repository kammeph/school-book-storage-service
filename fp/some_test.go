package fp_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	ints := []int{1, 2, 5, 1, 7, 1}
	assert.True(t, fp.Some(ints, func(i int) bool { return i == 1 }))
	assert.False(t, fp.Some(ints, func(i int) bool { return i == 8 }))

	strings := []string{"Hello", "World", "Hello"}
	assert.True(t, fp.Some(strings, func(s string) bool { return s == "Hello" }))
	assert.False(t, fp.Some(strings, func(s string) bool { return s == "You" }))
}
