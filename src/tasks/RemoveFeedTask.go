package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/gogap/logrus"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/feed"
	"github.com/inclee/gfeeds/src/storage/redis"
)

type RemoveFeedTask struct {
	activities []*activity.BaseActivty
	feed *feed.RedisFeed
}

func (t *RemoveFeedTask)ParseKwargs(kwargs map[string]interface{}) error {
	if act ,ok := kwargs["activities"];ok{
		ais := act.([]interface{})
		t.activities = make([]*activity.BaseActivty,len(ais),len(ais))
		for idx,ai := range ais{
			t.activities[idx] = new(activity.BaseActivty)
			if t.activities[idx] == nil {
				logrus.Error("create activity error")
			}
			err := json.Unmarshal([]byte(ai.(string)),t.activities[idx])
			if err != nil || t.activities[idx] == nil {
				logrus.Error("unmarshal activity error:",err.Error())
			}
		}
	}
	if user,ok := kwargs["user"];ok{
		t.feed = feed.NewRedisFeed()
		intu := int(user.(float64))
		logrus.Info("---> feed_user:",intu)
		t.feed.Init(intu,fmt.Sprint("in_feed_",intu),redis.NewRedisTimeLineStorage(new(redis.RedisTimeLineStorageDelegate)),&redis.ActiveStorage{})
	}
	return nil
}
func (t *RemoveFeedTask) RunTask() (interface{}, error) {
	t.feed.RemoveMany(t.activities)
	return nil, nil
}