package storage

import (
	"fmt"
	"github.com/inclee/gfeeds/activity"
	"hash/fnv"
)

type TimeLineStorager interface {
	Add(key string ,activity *activity.BaseActivty )
	AddMany(key string ,activties []*activity.BaseActivty ) int64
	RemoveMany(key string ,activties []*activity.BaseActivty ) int64
	GetActivities(key string,pgx int64,pgl int64)[]*activity.BaseActivty
}

type ActiveStorager interface {
	Add(key string ,activity *activity.BaseActivty )
	AddMany(key string ,activties []*activity.BaseActivty ) int64
}

type StoragerDelegate interface {
	AddToStorage(key string ,values[]*activity.BaseActivty)int64
	RemoveFromStorage(key string ,values[]*activity.BaseActivty)int64
	GetActivities(key string,pgx int64,pgl int64)[]*activity.BaseActivty
}

type BaseStorage struct {
	KeysNum int
}

func (self *BaseStorage)HashKey(key string)string{
	h := fnv.New32a()
	h.Write([]byte(key))
	return fmt.Sprintf("%s_key%d",int(h.Sum32()) % self.KeysNum)
}

type ActiveStorage struct {
	delegate StoragerDelegate
}

type TimeLineStorage struct {
	Delegate StoragerDelegate `json:"delegate"`
}

func NewTimeLineStorage(delegate StoragerDelegate) *TimeLineStorage{
	return &TimeLineStorage{
		Delegate:delegate,
	}
}

func(self *TimeLineStorage)Add(key string ,act*activity.BaseActivty ){
	self.AddMany(key,[]*activity.BaseActivty{ act})
}
func(self *TimeLineStorage)AddMany(key string ,activties []*activity.BaseActivty ) int64 {
	return self.Delegate.AddToStorage(key,activties)
}

func(self *TimeLineStorage)RemoveMany(key string ,activties []*activity.BaseActivty ) int64 {
	return self.Delegate.RemoveFromStorage(key,activties)
}

func(self *TimeLineStorage)GetActivities(key string,pgx int64,pgl int64)[]*activity.BaseActivty{
	return self.Delegate.GetActivities(key,pgx,pgl)
}

func(self *ActiveStorage)Add(key string ,act*activity.BaseActivty )()  {
	self.AddMany(key,[]*activity.BaseActivty{ act})
}
func(self *ActiveStorage)AddMany(key string ,activties []*activity.BaseActivty ) int64 {
	return self.delegate.AddToStorage(key,activties)
}


