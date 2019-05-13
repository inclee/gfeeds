package redis

import "github.com/gomodule/redigo/redis"

type RedisCache struct {
	key string
	pool *redis.Pool
	server string
}

func (self *RedisCache)Init(pool *redis.Pool)  {
	self.pool = pool
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

func (self RedisSortedSetCache)AddManay(key string,scores[]int,values[] interface{})(rcvn int, err error){
	c := self.pool.Get()
	defer  c.Close()
	params := make([]interface{},2*len(scores)+1,2*len(scores)+1)
	params[0] = key
	for idx,score := range  scores{
		params[2*idx+1] = score
		params[2*idx+2] = values[idx]
	}
	rcvn, err = redis.Int(c.Do("zadd",params...))
	c.Do("zrange feed_1s 0 -1")
	return
}

type RedisTimeLineCache struct {
	sortedSetCache *RedisSortedSetCache
	listCache *RedisListCache
	hashCache *RedisHaseCache
}

func (self *RedisTimeLineCache)Init(server string)  {
	pool := &redis.Pool{
		MaxIdle:16,
		MaxActive:0,
		IdleTimeout:300,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp",server)
		},
	}
	self.sortedSetCache = new(RedisSortedSetCache)
	self.listCache = new(RedisListCache)
	self.hashCache = new(RedisHaseCache)

	self.sortedSetCache.Init(pool)
	self.listCache.Init(pool)
	self.hashCache.Init(pool)
}