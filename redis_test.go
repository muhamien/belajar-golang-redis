package belajar_golang_redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func TestConnection(t *testing.T) {
	assert.NotNil(t, client)

	//err := client.Close()
	//assert.Nil(t, err)
}

var ctx = context.Background()

func TestPing(t *testing.T) {
	result, err := client.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.NotNil(t, "PONG", result)
}

func TestString(t *testing.T) {
	client.Set(ctx, "name", "Muhammad Amien", time.Second*3)
	result, err := client.Get(ctx, "name").Result()

	assert.Nil(t, err)
	assert.Equal(t, "Muhammad Amien", result)
}

func TestList(t *testing.T) {
	client.RPush(ctx, "names", "Muhammad")
	client.RPush(ctx, "names", "Amien")
	client.RPush(ctx, "names", "Rauf")

	assert.Equal(t, "Muhammad", client.LPop(ctx, "names").Val())
	assert.Equal(t, "Amien", client.LPop(ctx, "names").Val())
	assert.Equal(t, "Rauf", client.LPop(ctx, "names").Val())

	client.Del(ctx, "name")
}

func TestSet(t *testing.T) {
	client.SAdd(ctx, "fruits", "Pen")
	client.SAdd(ctx, "fruits", "Pen")
	client.SAdd(ctx, "fruits", "Pen")
	client.SAdd(ctx, "fruits", "Apple")
	client.SAdd(ctx, "fruits", "Apple")

	assert.Equal(t, int64(2), client.SCard(ctx, "fruits").Val())
	assert.Equal(t, []string{"Pen", "Apple"}, client.SMembers(ctx, "fruits").Val())
}
