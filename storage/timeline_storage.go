package storage

import (
	"encoding/json"
	"strconv"

	log "github.com/gogap/logrus"
	"github.com/inclee/gfeeds/activity"
)

var Cache *RedisTimeLineCache

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

func (self *RedisTimeLineStorageDelegate) AddToStorage(key string, activties []*activity.BaseActivty) int64 {
	scores := make([]float64, len(activties), len(activties))
	values := make([]interface{}, len(activties), len(activties))
	//ext := make([]string,len(activties),len(activties))
	for idx, act := range activties {
		if act == nil {
			log.Error("activity is nil")
		}
		if score, err := strconv.ParseFloat(act.SerializeId(),64); err == nil {
			scores[idx] = score
			values[idx], _ = act.Serialize()
		}
	}
	n, err := Cache.sortedSetCache.AddManay(key, scores, values)
	if err != nil {
		println(err.Error())
	}
	return n
}
func (self *RedisTimeLineStorageDelegate) RemoveFromStorage(key string, activties []*activity.BaseActivty) int64 {
	values := make([]interface{}, len(activties), len(activties))
	for idx, a := range activties {
		values[idx], _ = a.Serialize()
	}
	//ext := make([]string,len(activties),len(activties))
	n, err := Cache.sortedSetCache.RemoveManay(key, values)
	if err != nil {
		println(err.Error())
	}
	return n
}
func (self *RedisTimeLineStorageDelegate) GetActivities(key string, pgx int64, pgl int64) []*activity.BaseActivty {
	items, err := Cache.sortedSetCache.GetMany(key, pgx, pgl)
	if err != nil {
		log.Error(err.Error())
	}
	acts := make([]*activity.BaseActivty, 0, 0)
	for _, item := range items {
		act := new(activity.BaseActivty)
		if err := json.Unmarshal([]byte(item), act); err != nil {
			log.Error(err.Error())
		} else {
			acts = append(acts, act)
		}
	}
	return acts
}
