package main

import (
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/feed"
	"github.com/inclee/gfeeds/src/feedmanager"
	"github.com/inclee/gfeeds/src/storage"
	"github.com/inclee/gfeeds/src/storage/redis"
	"time"
)

type LZFeedManager struct {
	feedmanager.Manager
}
var manager *LZFeedManager
func init() {
	manager = new(LZFeedManager)
	if err := manager.Init();err !=nil{
		panic(err)
	}
}

func (m *LZFeedManager)getUserFeedClass()feed.Feed{
	return &feed.RedisFeed{&feed.BaseFeed{}}
}
func (m *LZFeedManager)getUserFollowerIds(uid int)([]int){
	return []int{1,2,3}
}
type UserFeed struct {
	feed.BaseFeed
}

func (self *UserFeed)KeyFormat()string{
	return "feed_%ds"
}

func (self *UserFeed)TimeLineStorage()storage.TimeLineStorage{
	return  &redis.RedisTimeLineStorage{}
}

//func (self *UserFeed)ActiveStorage()storage.ActiveStorage{
//	//return "feed_%ds"
//}

func main() {
	println("debug main")
	act := activity.NewActivity(
		"user:1",
		1,
		"post:1",
		"",
		time.Now(),
		"",
	)
	manager.AddActivity(1,act)
}