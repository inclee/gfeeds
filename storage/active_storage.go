package storage

import (
	"strconv"

	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/config"
)

type RedisActiveStorage struct {
	ActiveStorage
}

func (self *RedisActiveStorage) getCache(key string) RedisTimeLineCache {
	cache := RedisTimeLineCache{}
	cache.Init(config.Config.TimelineStorage)
	return cache
}
func (self *RedisActiveStorage) addToStorage(key string, activties []*activity.BaseActivty) (cnt int64, err error) {
	cache := self.getCache(key)
	scores := make([]float64, len(activties), len(activties))
	values := make([]interface{}, len(activties), len(activties))
	for idx, act := range activties {
		if score, err := strconv.ParseFloat(act.SerializeId(), 64); err == nil {
			scores[idx] = score
			values[idx] = act
		}
	}
	return cache.sortedSetCache.AddManay(key, scores, values)
}
