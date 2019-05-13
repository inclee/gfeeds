package main

import (
	"github.com/gocelery/gocelery"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/feed"
	"github.com/inclee/gfeeds/src/feedmanager"
	"github.com/inclee/gfeeds/src/storage"
	"github.com/inclee/gfeeds/src/storage/redis"
	"github.com/inclee/gfeeds/src/tasks"
	"time"
)

type FeedManagerDelegate   struct {
	feedmanager.Manager
}
var manager *feedmanager.Manager
func init() {
	manager = feedmanager.NewFeedManager(new(FeedManagerDelegate))
}
func (m *FeedManagerDelegate)GetPersonalFeed(uid int)feed.Feed{
	feed := &feed.RedisFeed{&feed.BaseFeed{}}
	feed.Init(uid,redis.DefaultRedisTimelineStorage,&redis.ActiveStorage{})
	return feed
}
func (m *FeedManagerDelegate)GetFollowIds(uid int)([]int){
	return []int{1,2,3,4,5,6,7,8,9}
}
type UserFeed struct {
	*feed.BaseFeed
}

func (self *UserFeed)KeyFormat()string{
	return "feed_%ds"
}

func (self *UserFeed)TimeLineStorage()storage.TimeLineStorager{
	return  &redis.RedisTimeLineStorage{}
}

func celery()  {
	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisCeleryBroker("redis://192.168.21.231:6379"),
		gocelery.NewRedisCeleryBackend("redis://192.168.21.231:6379"),
		5, // number of workers
	)
	cli.Register("feedmanager.add_activities_operation",&tasks.AddFeedTask{})
	cli.StartWorker()
	time.Sleep(60*time.Second)
}
func server()  {
		act := activity.NewActivity(
		"user:1",
		1,
		"post:1",
		"",
		time.Now(),
		`{
    "id": "6856bd64-f89f-4c4e-9eee-e82741b9fffe",
    "task": "feedmanager.add_activities_operation",}`,
	)
	manager.AddActivity(1,act)
}
func main() {
	println("debug main")
	//server()
	celery()
}