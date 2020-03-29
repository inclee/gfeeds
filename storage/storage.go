package storage

import (
	"fmt"
	"hash/fnv"

	"github.com/inclee/gfeeds/activity"
)

type TimeLineStorager interface {
	Add(key string, activity *activity.BaseActivty) error
	AddMany(key string, activties []*activity.BaseActivty) (int64, error)
	RemoveMany(key string, activties []*activity.BaseActivty) (int64, error)
	GetActivities(key string, pgx int64, pgl int64) ([]*activity.BaseActivty, error)
}

type ActiveStorager interface {
	Add(key string, activity *activity.BaseActivty) error
	AddMany(key string, activties []*activity.BaseActivty) (int64, error)
}

type StoragerDelegate interface {
	AddToStorage(key string, values []*activity.BaseActivty) (int64, error)
	RemoveFromStorage(key string, values []*activity.BaseActivty) (int64, error)
	GetActivities(key string, pgx int64, pgl int64) ([]*activity.BaseActivty, error)
}

type BaseStorage struct {
	KeysNum int
}

func (self *BaseStorage) HashKey(key string) string {
	h := fnv.New32a()
	h.Write([]byte(key))
	return fmt.Sprintf("%s_key%d", int(h.Sum32())%self.KeysNum)
}

type ActiveStorage struct {
	delegate StoragerDelegate
}

type TimeLineStorage struct {
	Delegate StoragerDelegate `json:"delegate"`
}

func NewTimeLineStorage(delegate StoragerDelegate) *TimeLineStorage {
	return &TimeLineStorage{
		Delegate: delegate,
	}
}

func (self *TimeLineStorage) Add(key string, act *activity.BaseActivty) (err error) {
	_, err = self.AddMany(key, []*activity.BaseActivty{act})
	return
}
func (self *TimeLineStorage) AddMany(key string, activties []*activity.BaseActivty) (cnt int64, err error) {
	return self.Delegate.AddToStorage(key, activties)
}

func (self *TimeLineStorage) RemoveMany(key string, activties []*activity.BaseActivty) (cnt int64, err error) {
	return self.Delegate.RemoveFromStorage(key, activties)
}

func (self *TimeLineStorage) GetActivities(key string, pgx int64, pgl int64) (acts []*activity.BaseActivty, err error) {
	return self.Delegate.GetActivities(key, pgx, pgl)
}

func (self *ActiveStorage) Add(key string, act *activity.BaseActivty) (err error) {
	_, err = self.AddMany(key, []*activity.BaseActivty{act})
	return
}
func (self *ActiveStorage) AddMany(key string, activties []*activity.BaseActivty) (cnt int64, err error) {
	return self.delegate.AddToStorage(key, activties)
}
