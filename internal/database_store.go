package internal

import (
	"context"
	_ "embed"
	"log"

	"github.com/go-redis/redis/v9"
)

//go:embed group.lua
var group string

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

// @todo output doesnt determine if theres too few users or room id is already used
// move 10 users to a new room with given id
func (s *DatabaseStore) Group(ctx context.Context, id string) (bool, error) {
	created, err := s.r.EvalSha(ctx, s.sha.group, []string{"lobby", "rooms", id}).Bool()
	if err != nil {
		log.Println(err)
		return false, err
	}

	return created, nil
}

func (s *DatabaseStore) RegisterGroupFunction(ctx context.Context) error {
	sha, err := s.r.ScriptLoad(ctx, group).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	s.sha.group = sha

	return nil
}
