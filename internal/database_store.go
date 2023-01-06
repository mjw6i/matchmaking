package internal

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
)

type DatabaseStore struct {
	r *redis.Client
}

func (s *DatabaseStore) Add(ctx context.Context, id string) error {
	_, err := s.r.SAdd(ctx, "lobby", id).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
