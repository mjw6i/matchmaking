package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	store := &DatabaseStore{}
	err := store.Add("123")
	assert.Nil(t, err)
}
