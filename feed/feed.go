package feed

import (
	"github.com/gogap/logrus"
	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/storage"
)

type Feed interface {
	Add(activty *activity.BaseActivty)
	AddMany(activties []*activity.BaseActivty) int64
	GetActivities(pgx int64, pgl int64) []*activity.BaseActivty
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

func (self *BaseFeed) Add(activty *activity.BaseActivty) {
	self.AddMany([]*activity.BaseActivty{activty})
}

func (self *BaseFeed) AddMany(activties []*activity.BaseActivty) int64 {
	var addCount int64
	if self.TimelineStorage != nil {
		addCount = self.TimelineStorage.AddMany(self.Key, activties)
	} else {
		logrus.Error("self.TimelineStorage")
	}
	return addCount
}

func (self *BaseFeed) RemoveMany(activties []*activity.BaseActivty) int64 {
	var addCount int64
	if self.TimelineStorage != nil {
		addCount = self.TimelineStorage.RemoveMany(self.Key, activties)
	} else {
		logrus.Error("self.TimelineStorage")
	}
	return addCount
}

func (self *BaseFeed) GetActivities(pgx int64, pgl int64) []*activity.BaseActivty {
	ret := make([]*activity.BaseActivty, 0, 0)
	if self.TimelineStorage != nil {
		ret = self.TimelineStorage.GetActivities(self.Key, pgx, pgl)
	} else {
		logrus.Error("self.TimelineStorage")
	}
	return ret
}
