package feed

import (
	"fmt"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/storage"
)

type Feed interface {
	KeyFormat()string
	Add(activty *activity.BaseActivty)
	AddMany(activties []*activity.BaseActivty)int
}

type BaseFeed struct {
	UserId int `json:"user_id"`
	Key string `json:"key"`
	TimelineStorage storage.TimeLineStorager `json:"timeline_storage"`
	ActiveStorage storage.ActiveStorager `json:"active_storage"`
}

func (self *BaseFeed)KeyFormat()string{
	return "feed_%ds"
}


func (self *BaseFeed)Init(userid int,timelineStorage storage.TimeLineStorager,activityStorage storage.ActiveStorager) {
	self.UserId= userid
	self.Key = fmt.Sprintf(self.KeyFormat(),userid)
	self.TimelineStorage = timelineStorage
	self.ActiveStorage =  activityStorage
}

func (self *BaseFeed)Add(activty *activity.BaseActivty){
	self.AddMany([]*activity.BaseActivty{activty})
}

func (self *BaseFeed)AddMany(activties []*activity.BaseActivty)int{
	addCount := self.TimelineStorage.AddMany(self.Key,activties)
	return addCount
}
