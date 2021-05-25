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

func (self *RedisTimeLineStorageDelegate) AddToStorage(key string, activties []*activity.BaseActivty) (cnt int64, err error) {
	scores := make([]float64, len(activties), len(activties))
	values := make([]interface{}, len(activties), len(activties))
	//ext := make([]string,len(activties),len(activties))
	for idx, act := range activties {
		if act == nil {
			log.Error("activity is nil")
		}
		if score, err := strconv.ParseFloat(act.SerializeId(), 64); err == nil {
			scores[idx] = score
			values[idx], _ = act.Serialize()
		}
	}
	cnt, err = Cache.sortedSetCache.AddManay(key, scores, values)
	return
}
func (self *RedisTimeLineStorageDelegate) RemoveFromStorage(key string, activties []*activity.BaseActivty) (cnt int64, err error) {
	values := make([]interface{}, len(activties), len(activties))
	for idx, a := range activties {
		values[idx], _ = a.Serialize()
	}
	//ext := make([]string,len(activties),len(activties))
	cnt, err = Cache.sortedSetCache.RemoveManay(key, values)
	return
}
func (self *RedisTimeLineStorageDelegate) GetActivities(key string, pgx int, pgl int) (acts []*activity.BaseActivty, err error) {
	items, err := Cache.sortedSetCache.GetMany(key, pgx, pgl)
	if err != nil {
		return
	}
	acts = make([]*activity.BaseActivty, 0, 0)
	for _, item := range items {
		act := new(activity.BaseActivty)
		if err := json.Unmarshal([]byte(item), act); err != nil {
			log.Error(err.Error())
		} else {
			acts = append(acts, act)
		}
	}
	return acts, err
}
