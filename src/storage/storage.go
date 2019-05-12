package storage

import (
	"fmt"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/define"
	"hash/fnv"
)

type TimeLineStorage interface {
	Add(key string ,activity *activity.BaseActivty )
	AddMany(key string ,activties []*activity.BaseActivty ) int16

}

type ActiveStorage interface {
	Add(key string ,activity *activity.BaseActivty )
	AddMany(key string ,activties []*activity.BaseActivty ) int16

}

type BaseStorage struct {
	KeysNum int
}

func (self *BaseStorage)HashKey(key string)string{
	h := fnv.New32a()
	h.Write([]byte(key))
	return fmt.Sprintf("%s_key%d",int(h.Sum32()) % self.KeysNum)
}

type BaseActiveStorage struct {

}

type BaseTimeLineStorage struct {

}

func(self *BaseTimeLineStorage)Add(key string ,act*activity.BaseActivty ){
	self.AddMany(key,[]*activity.BaseActivty{ act})
}
func(self *BaseTimeLineStorage)AddMany(key string ,activties []*activity.BaseActivty ) int16 {
	return self.addToStorage(key,activties)
}
func (self *BaseTimeLineStorage)addToStorage(key string ,activties []*activity.BaseActivty)  int16{
	panic(define.NotImplementedException)
}

func(self *BaseActiveStorage)Add(key string ,act*activity.BaseActivty )()  {
	self.AddMany(key,[]*activity.BaseActivty{ act})
}
func(self *BaseActiveStorage)AddMany(key string ,activties []*activity.BaseActivty ) int16 {
	return self.addToStorage(key,activties)
}
func (self *BaseActiveStorage)addToStorage(key string ,activties []*activity.BaseActivty)  int16{
	panic(define.NotImplementedException)
}