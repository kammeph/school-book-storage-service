package utils_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Setenv("ENV", "myvariable")
	v := utils.GetenvOrFallback("ENV", "fallback")
	assert.Equal(t, v, "myvariable")
}

func TestGetEnvFallback(t *testing.T) {
	v := utils.GetenvOrFallback("ENV", "fallback")
	assert.Equal(t, v, "fallback")
}
