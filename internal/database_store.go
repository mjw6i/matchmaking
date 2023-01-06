package internal

import "github.com/go-redis/redis/v9"

type DatabaseStore struct {
	r *redis.Client
}

func (s *DatabaseStore) Add(id string) error {
	return nil
}
