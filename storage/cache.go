package storage

import (
	"github.com/go-redis/redis"
)

type RedisCache struct {
	key     string
	cluster *redis.ClusterClient
	server  string
}

func (self *RedisCache) Init(cluster *redis.ClusterClient) {
	self.cluster = cluster
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
	rcvn, err = self.cluster.ZAdd(key, params...).Result()
	return
}

func (self RedisSortedSetCache) RemoveManay(key string, values []interface{}) (rcvn int64, err error) {
	rcvn, err = self.cluster.ZRem(key, values...).Result()
	return
}

func (self RedisSortedSetCache) GetMany(key string, pgx int64, pgl int64) ([]string, error) {
	rets, err := self.cluster.ZRevRangeByScore(key, redis.ZRangeBy{
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
	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    servers,
		Password: "v7baLZM9efsgvRI7",
	})
	self.sortedSetCache = new(RedisSortedSetCache)
	self.listCache = new(RedisListCache)
	self.hashCache = new(RedisHaseCache)

	self.sortedSetCache.Init(cluster)
	self.listCache.Init(cluster)
	self.hashCache.Init(cluster)
}
