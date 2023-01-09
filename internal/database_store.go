package internal

import (
	"context"
	"errors"
	"log"

	"github.com/go-redis/redis/v9"
)

type DatabaseStore struct {
	r   *redis.Client
	sha struct { //@todo Most likely will be moved to a redis-managed function
		group string // @todo should that be initialised by New()?
	}
}

func (s *DatabaseStore) Add(ctx context.Context, id string, score float64) error {
	_, err := s.r.ZAdd(ctx, "lobby", redis.Z{Score: score, Member: id}).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *DatabaseStore) Group(ctx context.Context) ([]string, error) {
	res, err := s.r.EvalSha(ctx, s.sha.group, []string{"lobby"}).StringSlice()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(res)%2 != 0 {
		return nil, errors.New("invalid length")
	}

	ids := make([]string, 0, len(res)/2)

	for i, e := range res {
		if i%2 == 0 {
			ids = append(ids, e)
		}
	}

	return ids, nil
}

func (s *DatabaseStore) RegisterGroupFunction(ctx context.Context) error {
	// @todo if my app crashes after getting the ids, users will be removed from the lobby but not put inside a room
	f := `
		local count = redis.call('ZCARD', KEYS[1])
		if count < 10 then
			return {}
		end

		return redis.call('ZPOPMIN', KEYS[1], 10)
	`

	sha, err := s.r.ScriptLoad(ctx, f).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	s.sha.group = sha

	return nil
}
