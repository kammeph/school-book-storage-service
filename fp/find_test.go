package fp_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/fp"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	ints := []int{1, 2, 5, 8, 7, 4}
	intResult := fp.Find(ints, func(i int) bool { return i == 5 })
	assert.Equal(t, 5, *intResult)

	intResult = fp.Find(ints, func(i int) bool { return i == 3 })
	assert.Nil(t, intResult)

	strings := []string{"Hello", "World"}
	stringResult := fp.Find(strings, func(s string) bool { return s == "Hello" })
	assert.Equal(t, "Hello", *stringResult)

	stringResult = fp.Find(strings, func(s string) bool { return s == "You" })
	assert.Nil(t, stringResult)
}

func TestFindIndex(t *testing.T) {
	ints := []int{1, 2, 5, 8, 7, 4}
	intIndex := fp.FindIndex(ints, func(i int) bool { return i == 5 })
	assert.Equal(t, 2, intIndex)

	intIndex = fp.FindIndex(ints, func(i int) bool { return i == 3 })
	assert.Equal(t, -1, intIndex)

	strings := []string{"Hello", "World"}
	stringIndex := fp.Find(strings, func(s string) bool { return s == "Hello" })
	assert.Equal(t, "Hello", *stringIndex)

	stringIndex = fp.Find(strings, func(s string) bool { return s == "You" })
	assert.Nil(t, stringIndex)
}
