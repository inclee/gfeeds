package redis

import (
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/storage"
	"strconv"
)

type RedisTimeLineStorage struct {
	*storage.BaseTimeLineStorage
}

func (self* RedisTimeLineStorage)getCache(key string)RedisTimeLineCache{
	cache := RedisTimeLineCache{}
	cache.Init("localhost:6379")
	return cache
}
func (self *RedisTimeLineStorage)addToStorage(key string ,activties []*activity.BaseActivty)  int{
	cache := self.getCache(key)
	scores := make([]int,len(activties),len(activties))
	values := make([]interface{},len(activties),len(activties))
	for idx,act := range activties{
		if score, err := strconv.Atoi(act.SerializeId()) ;err == nil {
			scores[idx] = score
			values[idx] = act
		}
	}
	cache.sortedSetCache.AddManay(key,scores,values)
	return len(scores)
}