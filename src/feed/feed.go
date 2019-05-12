package feed

import (
	"fmt"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/storage"
)

type Feed interface {
	KeyFormat()string
	Add(activty *activity.BaseActivty)
	AddMany(activties []*activity.BaseActivty)int16
	TimeLineStorage() storage.TimeLineStorage
	ActiveStorage() storage.ActiveStorage
}

type BaseFeed struct {
	userId int
	key string
	timelineStorage storage.TimeLineStorage
	activeStorage storage.ActiveStorage
}

func (self *BaseFeed)KeyFormat()string{
	return "feed_%ds"
}

func (self *BaseFeed)TimeLineStorage()storage.TimeLineStorage{
	return  &storage.BaseTimeLineStorage{}
}

func (self *BaseFeed)ActiveStorage()storage.ActiveStorage{
	return &storage.BaseActiveStorage{}
}



func (self *BaseFeed)Init(userid int) {
	self.userId = userid
	self.key = fmt.Sprintf(self.KeyFormat(),userid)
	self.timelineStorage = self.TimeLineStorage()
	self.activeStorage = self.ActiveStorage()
}

func (self *BaseFeed)Add(activty *activity.BaseActivty){
	self.AddMany([]*activity.BaseActivty{activty})
}

func (self *BaseFeed)AddMany(activties []*activity.BaseActivty)int16{
	addCount := self.timelineStorage.AddMany(self.key,activties)
	return addCount
}
