package internal

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	r := getRedis()
	flushRedis(t, r)
	store := getStore(r)
	err := store.Add(context.Background(), "123")
	require.Nil(t, err)
	requireSet(t, r, "lobby", []string{"123"})
}

func requireSet(t *testing.T, r *redis.Client, name string, expected []string) {
	actual, err := r.SMembers(context.Background(), name).Result()
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func getStore(r *redis.Client) *DatabaseStore {
	return &DatabaseStore{r: r}
}

func flushRedis(t *testing.T, r *redis.Client) {
	err := r.FlushAll(context.Background()).Err()
	require.Nil(t, err)
}

func getRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}
