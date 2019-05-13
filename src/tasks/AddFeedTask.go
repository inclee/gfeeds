package tasks

import (
	"encoding/json"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/feed"
	"github.com/inclee/gfeeds/src/storage/redis"
)

type AddFeedTask struct {
	activities []*activity.BaseActivty
	feed *feed.RedisFeed
}

func (t *AddFeedTask)ParseKwargs(kwargs map[string]interface{}) error {
	if act ,ok := kwargs["activities"];ok{
		ais := act.([]interface{})
		t.activities = make([]*activity.BaseActivty,len(ais),len(ais))
		for idx,ai := range ais{
			t.activities[idx] = new(activity.BaseActivty)
			println(ai.(string))
			json.Unmarshal([]byte(ai.(string)),t.activities[idx])
		}
	}
	if user,ok := kwargs["user"];ok{
		t.feed = feed.NewRedisFeed()
		t.feed.Init(int(user.(float64)),redis.NewRedisTimeLineStorage(new(redis.RedisTimeLineStorageDelegate)),&redis.ActiveStorage{})
	}
	return nil
}
func (t *AddFeedTask) RunTask() (interface{}, error) {
	t.feed.AddMany(t.activities)
	return nil, nil
}