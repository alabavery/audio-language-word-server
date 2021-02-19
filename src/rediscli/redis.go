package rediscli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

var ctx = context.TODO()

// WordRedisCli is a wrapper for the client
type WordRedisCli struct {
	rdb *redis.Client
}

// GetWordRedisCli gives WordRedisCli
func GetWordRedisCli() *WordRedisCli {
	var cli WordRedisCli
	cli.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &cli
}

// Get gets key in the parts of speech
func (c *WordRedisCli) Get(key string) ([]byte, bool) {
	val, err := c.rdb.Get(key).Result()
	if err == redis.Nil {
		return json.RawMessage(nil), false
	} else if err != nil {
		panic(err)
	}
	return json.RawMessage(val), true
}

func (c *WordRedisCli) Set(key string, val []byte) {
	err := c.rdb.Set(key, val, 0).Err()
	if err != nil {
		panic(err)
	}
}

func WordKeyFromId(id int) string {
	return fmt.Sprint(id)
}
