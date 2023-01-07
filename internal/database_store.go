package internal

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
)

type DatabaseStore struct {
	r *redis.Client
}

func (s *DatabaseStore) Add(ctx context.Context, id string, score float64) error {
	_, err := s.r.ZAdd(ctx, "lobby", redis.Z{Score: score, Member: id}).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *DatabaseStore) Group(ctx context.Context) error {
	return nil
}
