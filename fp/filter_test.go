package fp_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	ints := []int{1, 2, 5, 1, 7, 1}
	intResults := fp.Filter(ints, func(i int) bool { return i != 1 })
	assert.EqualValues(t, []int{2, 5, 7}, intResults)

	strings := []string{"Hello", "World", "Hello"}
	stringResults := fp.Filter(strings, func(s string) bool { return s == "Hello" })
	assert.EqualValues(t, []string{"Hello", "Hello"}, stringResults)
}
