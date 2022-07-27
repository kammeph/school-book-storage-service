package fp_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	ints := []int{1, 2, 5, 1, 7, 1}
	ints = fp.Remove(ints, func(i int) bool { return i == 7 })
	assert.EqualValues(t, []int{1, 2, 5, 1, 1}, ints)

	ints = fp.Remove(ints, func(i int) bool { return i == 8 })
	assert.EqualValues(t, []int{1, 2, 5, 1, 1}, ints)

	strings := []string{"Hello", "World", "Hello"}
	strings = fp.Remove(strings, func(s string) bool { return s == "World" })
	assert.EqualValues(t, []string{"Hello", "Hello"}, strings)

	strings = fp.Remove(strings, func(s string) bool { return s == "You" })
	assert.EqualValues(t, []string{"Hello", "Hello"}, strings)
}

func TestRemoveIndex(t *testing.T) {
	ints := []int{1, 2, 5, 1, 7, 1}
	ints = fp.RemoveIndex(ints, 4)
	assert.EqualValues(t, []int{1, 2, 5, 1, 1}, ints)

	ints = fp.RemoveIndex(ints, -1)
	assert.EqualValues(t, []int{1, 2, 5, 1, 1}, ints)

	strings := []string{"Hello", "World", "Hello"}
	strings = fp.RemoveIndex(strings, 1)
	assert.EqualValues(t, []string{"Hello", "Hello"}, strings)

	strings = fp.RemoveIndex(strings, -1)
	assert.EqualValues(t, []string{"Hello", "Hello"}, strings)
}
