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

func (self RedisSortedSetCache)AddManay(key string,scores[]int,values[] interface{}){
	c := self.pool.Get()
	defer  c.Close()
	params := make([]interface{},2*len(scores)+1,2*len(scores)+1)
	params = append(params, key)
	for idx,score := range  scores{
		params = append(params, score)
		params = append(params, values[idx])
	}
	c.Do("zadd",params...)
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
	self.sortedSetCache.Init(pool)
	self.listCache.Init(pool)
	self.hashCache.Init(pool)
}