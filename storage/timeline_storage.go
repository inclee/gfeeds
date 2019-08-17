package storage

import (
	"encoding/json"
	"hash/fnv"
	"strconv"

	log "github.com/gogap/logrus"
	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/config"
)

var RedisTimeLineCaches []*RedisTimeLineCache

type RedisTimeLineStorage struct {
	*TimeLineStorage
}

var DefaultRedisTimelineStorage = NewRedisTimeLineStorage(new(RedisTimeLineStorageDelegate))

func NewRedisTimeLineStorage(delegate StoragerDelegate) *RedisTimeLineStorage {
	rs := new(RedisTimeLineStorage)
	rs.TimeLineStorage = &TimeLineStorage{
		Delegate: delegate,
	}
	return rs
}

type RedisTimeLineStorageDelegate struct {
}

func (self *RedisTimeLineStorageDelegate) getCache(key string) *RedisTimeLineCache {
	h := fnv.New32a()
	h.Write([]byte(key))
	idx := int(h.Sum32()) % len(config.Config.TimelineStorage)
	log.Info("select redis :", config.Config.TimelineStorage[idx])
	return RedisTimeLineCaches[idx]
}
func (self *RedisTimeLineStorageDelegate) AddToStorage(key string, activties []*activity.BaseActivty) int {
	cache := self.getCache(key)
	scores := make([]int, len(activties), len(activties))
	values := make([]interface{}, len(activties), len(activties))
	//ext := make([]string,len(activties),len(activties))
	for idx, act := range activties {
		if act == nil {
			log.Error("activity is nil")
		}
		if score, err := strconv.Atoi(act.SerializeId()); err == nil {
			scores[idx] = score
			values[idx], _ = act.Serialize()
		}
	}
	n, err := cache.sortedSetCache.AddManay(key, scores, values)
	if err != nil {
		println(err.Error())
	}
	return n
}
func (self *RedisTimeLineStorageDelegate) RemoveFromStorage(key string, activties []*activity.BaseActivty) int {
	cache := self.getCache(key)
	values := make([]interface{}, len(activties), len(activties))
	for idx, a := range activties {
		values[idx], _ = a.Serialize()
	}
	//ext := make([]string,len(activties),len(activties))
	n, err := cache.sortedSetCache.RemoveManay(key, values)
	if err != nil {
		println(err.Error())
	}
	return n
}
func (self *RedisTimeLineStorageDelegate) GetActivities(key string, pgx int, pgl int) []*activity.BaseActivty {
	cache := self.getCache(key)
	items, err := cache.sortedSetCache.GetMany(key, pgx, pgl)
	if err != nil {
		log.Error(err.Error())
	}
	acts := make([]*activity.BaseActivty, 0, 0)
	for _, item := range items {
		if bits, ok := item.([]byte); ok {
			act := new(activity.BaseActivty)
			if err := json.Unmarshal(bits, act); err != nil {
				log.Error(err.Error())
			} else {
				acts = append(acts, act)
			}
		} else {
			log.Error("cover type failed", item)
		}
	}
	return acts
}
