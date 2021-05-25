package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/gogap/logrus"
	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/feed"
	"github.com/inclee/gfeeds/storage"
)

type RemoveFeedTask struct {
	activities []*activity.BaseActivty
	feed       *feed.BaseFeed
}

func (t *RemoveFeedTask) ParseKwargs(kwargs map[string]interface{}) error {
	if act, ok := kwargs["activities"]; ok {
		ais := act.([]interface{})
		t.activities = make([]*activity.BaseActivty, len(ais), len(ais))
		for idx, ai := range ais {
			t.activities[idx] = new(activity.BaseActivty)
			if t.activities[idx] == nil {
				logrus.Error("create activity error")
			}
			err := json.Unmarshal([]byte(ai.(string)), t.activities[idx])
			if err != nil || t.activities[idx] == nil {
				logrus.Error("unmarshal activity error:", err.Error())
			}
		}
	}
	if user, ok := kwargs["user"]; ok {
		t.feed = &feed.BaseFeed{}
		intu := uint64(user.(float64))
		logrus.Info("---> feed_user:", intu)
		t.feed.Init(intu, fmt.Sprintf(kwargs["target_feed_id"].(string), intu), storage.NewRedisTimeLineStorage(new(storage.RedisTimeLineStorageDelegate)), &storage.ActiveStorage{})
	}
	return nil
}
func (t *RemoveFeedTask) RunTask() (interface{}, error) {
	t.feed.RemoveMany(t.activities)
	return nil, nil
}
