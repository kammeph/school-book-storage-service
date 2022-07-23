package utils_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/infrastructure/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		setEnv   bool
		env      string
		fallback string
		expected string
	}{
		{
			name:     "env is set",
			setEnv:   true,
			env:      "my ENV",
			fallback: "my ENV fallback",
			expected: "my ENV",
		},
		{
			name:     "env is not set",
			setEnv:   false,
			env:      "my ENV",
			fallback: "my ENV fallback",
			expected: "my ENV fallback",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setEnv {
				t.Setenv("ENV", test.env)
			}
			v := utils.GetenvOrFallback("ENV", test.fallback)
			assert.Equal(t, test.expected, v)
		})
	}
}
