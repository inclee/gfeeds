package storage

import (
	"fmt"
	"github.com/inclee/gfeeds/src/activity"
	"hash/fnv"
)

type TimeLineStorager interface {
	Add(key string ,activity *activity.BaseActivty )
	AddMany(key string ,activties []*activity.BaseActivty ) int
}

type ActiveStorager interface {
	Add(key string ,activity *activity.BaseActivty )
	AddMany(key string ,activties []*activity.BaseActivty ) int
}

type StoragerDelegate interface {
	AddToStorage(key string ,values[]*activity.BaseActivty)int
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
func(self *TimeLineStorage)AddMany(key string ,activties []*activity.BaseActivty ) int {

	return self.Delegate.AddToStorage(key,activties)
}

func(self *ActiveStorage)Add(key string ,act*activity.BaseActivty )()  {
	self.AddMany(key,[]*activity.BaseActivty{ act})
}
func(self *ActiveStorage)AddMany(key string ,activties []*activity.BaseActivty ) int {
	return self.delegate.AddToStorage(key,activties)
}
