package redis

import (
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/storage"
	"strconv"
)

type RedisTimeLineStorage struct {
	*storage.TimeLineStorage
}

var DefaultRedisTimelineStorage = NewRedisTimeLineStorage( new(RedisTimeLineStorageDelegate))

func NewRedisTimeLineStorage(delegate storage.StoragerDelegate)*RedisTimeLineStorage {
	rs := new (RedisTimeLineStorage)
	rs.TimeLineStorage = &storage.TimeLineStorage{
		Delegate:delegate,
	}
	return rs
}

type RedisTimeLineStorageDelegate struct {

}
func (self* RedisTimeLineStorageDelegate)getCache(key string)RedisTimeLineCache{
	cache := RedisTimeLineCache{}
	cache.Init("192.168.21.231:6379")
	return cache
}
func (self *RedisTimeLineStorageDelegate)AddToStorage(key string ,activties []*activity.BaseActivty)  int{
	cache := self.getCache(key)
	scores := make([]int,len(activties),len(activties))
	values := make([]interface{},len(activties),len(activties))
	for idx,act := range activties{
		if score, err := strconv.Atoi(act.SerializeId()) ;err == nil {
			scores[idx] = score
			values[idx],_ = act.JsonSerialize()
		}
	}
	n,err := cache.sortedSetCache.AddManay(key,scores,values)
	if err != nil{
		println(err.Error())
	}
	return n
}