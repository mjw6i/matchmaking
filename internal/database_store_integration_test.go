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

	res, err := store.Group(context.Background(), "test")
	require.Nil(t, err)
	require.False(t, res)

	requireSet(t, r, "lobby", []string{"1a", "1b", "1c"})
}

func TestGroupCreatesOnlyFreshRooms(t *testing.T) {
	r := getRedis()
	flushRedis(t, r)
	store := getStore(r)
	err := store.RegisterGroupFunction(context.Background())
	require.Nil(t, err)
	micro := time.Now().UnixMicro()
	users1 := []string{
		"1a", "1b", "1c", "1d", "1e", "1f", "1g", "1h", "1i", "1j",
	}
	users2 := []string{
		"2a", "2b", "2c", "2d", "2e", "2f", "2g", "2h", "2i", "2j",
	}
	for _, u := range users1 {
		err := store.Add(context.Background(), u, float64(micro))
		require.Nil(t, err)
		micro++
	}
	for _, u := range users2 {
		err := store.Add(context.Background(), u, float64(micro))
		require.Nil(t, err)
		micro++
	}

	res, err := store.Group(context.Background(), "test")
	require.Nil(t, err)
	require.True(t, res)

	res, err = store.Group(context.Background(), "test")
	require.Nil(t, err)
	require.False(t, res)

	requireSet(t, r, "lobby", users2)
}

func TestGroupOrderIsDeterminedByJoinTime(t *testing.T) {
	r := getRedis()
	flushRedis(t, r)
	store := getStore(r)
	err := store.RegisterGroupFunction(context.Background())
	require.Nil(t, err)

	micro := time.Now().UnixMicro()
	users := []string{
		"1a", "1b", "1c", "1d", "1e", "1f", "1g",
		"3a", "3b", "3c",
		"2a", "2b", "2c",
	}
	for _, u := range users {
		err := store.Add(context.Background(), u, float64(micro))
		require.Nil(t, err)
		micro++
	}

	_, err = store.Group(context.Background(), "test")
	require.Nil(t, err)

	requireSet(t, r, "lobby", []string{"2a", "2b", "2c"})
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
