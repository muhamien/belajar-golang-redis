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

func TestSortedSet(t *testing.T) {
	client.ZAdd(ctx, "scores", redis.Z{Score: 100, Member: "Muhammad"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 85, Member: "Amien"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 96, Member: "Rauf"})

	assert.Equal(t, []string{"Amien", "Rauf", "Muhammad"}, client.ZRange(ctx, "scores", 0, 2).Val())
	assert.Equal(t, "Muhammad", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Rauf", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Amien", client.ZPopMax(ctx, "scores").Val()[0].Member)

}

func TestHash(t *testing.T) {
	client.HSet(ctx, "user:1", "id", "1", "name", "Amien")
	client.HSet(ctx, "user:2", "id", "1", "name", "Rauf")

	user1 := client.HGetAll(ctx, "user:1").Val()
	user2 := client.HGetAll(ctx, "user:2").Val()

	assert.Equal(t, "1", user1["id"])
	assert.Equal(t, "Amien", user1["name"])

	assert.Equal(t, "1", user2["id"])
	assert.Equal(t, "Rauf", user2["name"])

	client.Del(ctx, "user:1")
}

func TestGeoPoint(t *testing.T) {
	client.GeoAdd(ctx, "sellers", &redis.GeoLocation{
		Name:      "Toko A",
		Longitude: 110.408456,
		Latitude:  -7.739936,
	})
	client.GeoAdd(ctx, "sellers", &redis.GeoLocation{
		Name:      "Toko B",
		Longitude: 110.4144109,
		Latitude:  -7.7420133,
	})

	distance := client.GeoDist(ctx, "sellers", "Toko A", "Toko B", "km").Val()
	assert.Equal(t, 0.6963, distance)

	sellers := client.GeoSearch(ctx, "sellers", &redis.GeoSearchQuery{
		Longitude:  110.410960,
		Latitude:   -7.740996,
		Radius:     15,
		RadiusUnit: "km",
	}).Val()

	assert.Equal(t, []string{"Toko A", "Toko B"}, sellers)
}

func TestHyperLogLog(t *testing.T) {
	client.PFAdd(ctx, "visitors", "Muhammad", "Amien", "Rauf")
	client.PFAdd(ctx, "visitors", "Joko", "Gibs", "Prabs")
	client.PFAdd(ctx, "visitors", "Gibs", "Amien", "Prabs")

	total := client.PFCount(ctx, "visitors").Val()
	assert.Equal(t, int64(6), total)
}

func TestPipeline(t *testing.T) {
	_, err := client.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SetEx(ctx, "name", "Eko", 5*time.Second)
		pipeliner.SetEx(ctx, "address", "Indonesia", 5*time.Second)
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, "Eko", client.Get(ctx, "name").Val())
	assert.Equal(t, "Indonesia", client.Get(ctx, "address").Val())
}

func TestTransaction(t *testing.T) {
	_, err := client.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SetEx(ctx, "name", "Joko", 5*time.Second)
		pipeliner.SetEx(ctx, "address", "Cirebon", 5*time.Second)
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, "Joko", client.Get(ctx, "name").Val())
	assert.Equal(t, "Cirebon", client.Get(ctx, "address").Val())
}
