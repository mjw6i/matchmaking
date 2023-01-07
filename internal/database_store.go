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
	// @todo for now only removes users
	err := s.r.Watch(ctx, transaction(ctx), "lobby")

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func transaction(ctx context.Context) func(tx *redis.Tx) error {
	return func(tx *redis.Tx) error {
		count, err := tx.ZCard(ctx, "lobby").Result()
		if err != nil {
			return err
		}

		if count < 10 {
			return nil
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err = pipe.ZPopMin(ctx, "lobby", 10).Result()
			return err
		})

		if err != nil {
			return err
		}

		return nil
	}
}
