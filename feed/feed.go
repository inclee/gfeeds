package feed

import (
	"fmt"

	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/storage"
)

type Feed interface {
	Add(activty *activity.BaseActivty) error
	AddMany(activties []*activity.BaseActivty) (int64, error)
	GetActivities(pgx int64, pgl int64) ([]*activity.BaseActivty, error)
}

type BaseFeed struct {
	UserId          int                      `json:"user_id"`
	Key             string                   `json:"key"`
	TimelineStorage storage.TimeLineStorager `json:"timeline_storage"`
	ActiveStorage   storage.ActiveStorager   `json:"active_storage"`
}

func (self *BaseFeed) Init(userid int, key string, timelineStorage storage.TimeLineStorager, activityStorage storage.ActiveStorager) {
	self.UserId = userid
	self.Key = key
	self.TimelineStorage = timelineStorage
	self.ActiveStorage = activityStorage
}

func (self *BaseFeed) Add(activty *activity.BaseActivty) (err error) {
	_, err = self.AddMany([]*activity.BaseActivty{activty})
	return
}

func (self *BaseFeed) AddMany(activties []*activity.BaseActivty) (addCount int64, err error) {
	if self.TimelineStorage != nil {
		addCount, err = self.TimelineStorage.AddMany(self.Key, activties)
	} else {
		err = fmt.Errorf("timeline storeage is null")
	}
	return
}

func (self *BaseFeed) RemoveMany(activties []*activity.BaseActivty) (cnt int64, err error) {
	if self.TimelineStorage != nil {
		cnt, err = self.TimelineStorage.RemoveMany(self.Key, activties)
	} else {
		err = fmt.Errorf("timeline storage is null")
	}
	return
}

func (self *BaseFeed) GetActivities(pgx int64, pgl int64) (acts []*activity.BaseActivty, err error) {
	if self.TimelineStorage != nil {
		acts, err = self.TimelineStorage.GetActivities(self.Key, pgx, pgl)
	} else {
		err = fmt.Errorf("time line storage is null")
	}
	return
}
