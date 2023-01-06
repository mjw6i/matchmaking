package internal

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	store := getStore()
	err := store.Add(context.Background(), "123")
	assert.Nil(t, err)
}

func getStore() *DatabaseStore {
	r := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	store := &DatabaseStore{r: r}
	return store
}
