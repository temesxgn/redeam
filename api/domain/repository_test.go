package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepo_NewRepositoryMissingEnvVariables(t *testing.T) {
	_, err := NewRepository()
	assert.NotNil(t, err)
}
