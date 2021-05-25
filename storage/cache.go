package storage

import (
	"github.com/go-redis/redis"
)

var RedisClent *redis.Client

type RedisCache struct {
	key    string
	conn   *redis.Client
	server string
}

func (self *RedisCache) Init(conn *redis.Client) {
	self.conn = conn
}

type RedisListCache struct {
	RedisCache
}

type RedisHaseCache struct {
	RedisCache
}

type RedisSortedSetCache struct {
	RedisCache
}

func (self RedisSortedSetCache) AddManay(key string, scores []float64, values []interface{}) (rcvn int64, err error) {
	params := make([]redis.Z, len(scores))
	for idx, score := range scores {
		params[idx] = redis.Z{
			Score:  score,
			Member: values[idx],
		}
	}
	rcvn, err = self.conn.ZAdd(key, params...).Result()
	return
}

func (self RedisSortedSetCache) RemoveManay(key string, values []interface{}) (rcvn int64, err error) {
	rcvn, err = self.conn.ZRem(key, values...).Result()
	return
}

func (self RedisSortedSetCache) GetMany(key string, pgx int, pgl int) ([]string, error) {
	rets, err := self.conn.ZRevRangeByScore(key, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(pgx * pgl),
		Count:  int64(pgl),
	},
	).Result()
	if err != nil {
		return nil, err
	} else {
		return rets, nil
	}
}

type RedisTimeLineCache struct {
	sortedSetCache *RedisSortedSetCache
	listCache      *RedisListCache
	hashCache      *RedisHaseCache
}

func (self *RedisTimeLineCache) Init(servers []string) {
	// conn := redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: servers,
	// 	// Password: "v7baLZM9efsgvRI7",
	// })
	conn := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RedisClent = conn
	self.sortedSetCache = new(RedisSortedSetCache)
	self.listCache = new(RedisListCache)
	self.hashCache = new(RedisHaseCache)

	self.sortedSetCache.Init(conn)
	self.listCache.Init(conn)
	self.hashCache.Init(conn)
}
