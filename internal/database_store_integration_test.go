package internal

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	r := getRedis()
	flushRedis(t, r)
	store := getStore(r)
	// @todo check if UnixMicro mapped to a float generates a sensible score
	micro := time.Now().UnixMicro()
	err := store.Add(context.Background(), "123", float64(micro))
	require.Nil(t, err)
	requireSet(t, r, "lobby", []string{"123"})
}

func TestGroupCreatesOnlyFullRooms(t *testing.T) {
	r := getRedis()
	flushRedis(t, r)
	store := getStore(r)
	err := store.RegisterGroupFunction(context.Background())
	require.Nil(t, err)
	micro := time.Now().UnixMicro()
	users := []string{"1a", "1b", "1c"}
	for _, u := range users {
		err := store.Add(context.Background(), u, float64(micro))
		require.Nil(t, err)
	}

	ids, err := store.Group(context.Background())
	require.Nil(t, err)
	require.Empty(t, make([]string, 0))
	require.Empty(t, ids)

	requireSet(t, r, "lobby", []string{"1a", "1b", "1c"})
}

func requireSet(t *testing.T, r *redis.Client, name string, expected []string) {
	actual, err := r.ZRange(context.Background(), name, 0, -1).Result()
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
