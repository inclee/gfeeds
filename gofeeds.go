package main

import (
	"time"
)

// type FeedManagerDelegate struct {
// 	feedmanager.Manager
// }

// var manager *feedmanager.Manager

// func init() {
// 	PrepareRedisTimeline()
// 	manager = feedmanager.NewFeedManager(new(FeedManagerDelegate))
// }
// func (m *FeedManagerDelegate) GetPersonalFeed(uid int) feed.Feed {
// 	feed := &feed.RedisFeed{&feed.BaseFeed{}}
// 	feed.Init(uid, fmt.Sprint("out_feed_", uid), storage.DefaultRedisTimelineStorage, &storage.ActiveStorage{})
// 	return feed
// }
// func (m *FeedManagerDelegate) GetFollowIds(uid int) []int {
// 	ret := make([]int, 0, 0)
// 	for i := 0; i < 50; i++ {
// 		ret = append(ret, rand.Int()%100000)
// 	}
// 	return ret
// }

// type UserFeed struct {
// 	*feed.RedisFeed
// }

// func (self *UserFeed) KeyFormat() string {
// 	return "feed_%ds"
// }

// func (self *UserFeed) TimeLineStorage() storage.TimeLineStorager {
// 	return &redis.RedisTimeLineStorage{}
// }

// func celery() {
// 	cli, _ := gocelery.NewCeleryClient(
// 		gocelery.NewRedisCeleryBroker(config.GConfig.Celery.Broker),
// 		gocelery.NewRedisCeleryBackend(config.GConfig.Celery.Backend),
// 		config.GConfig.Celery.WorkNum, // number of workers
// 	)
// 	cli.Register("feedmanager.add_activities_operation", &tasks.AddFeedTask{})
// 	cli.StartWorker()

// }
// func server(user int) {
// 	act := activity.NewActivity(
// 		fmt.Sprintf("user:%d", user),
// 		1,
// 		fmt.Sprintf("post:%d", user),
// 		"post:1",
// 		time.Now(),
// 		`{
//    "id": "6856bd64-f89f-4c4e-9eee-e82741b9fffe",
// 	"content:Imagine I have a document written by somebody which contains plain text. Now that plain text just happens to be valid JSON as well. Would I then be wrong to use text/plain as its mime-type? JSON is a SUB-TYPE of text. So I think both should be allowed. The question is which works better in practice. According to comment by codetoshare IE has problems with application/json. But no browser should have problems with text/plain. If text/plain is unsafe then how can I serve text-files from my web-site? "
//    "task": "feedmanager.add_activities_operation",}`,
// 	)
// 	manager.AddActivity(user, act)
// }
// func PrepareRedisTimeline() {
// 	for _, r := range config.GConfig.TimelineStorage.Redis {
// 		cache := new(redis.RedisTimeLineCache)
// 		cache.Init(r)
// 		redis.RedisTimeLineCaches = append(redis.RedisTimeLineCaches, cache)
// 	}
// }
// func getPersonnalFeed(user int) {
// 	acts := manager.LoadPersonFeeds(user, 1, 10)
// 	for _, act := range acts {
// 		b, _ := json.Marshal(act)
// 		logrus.Info("personal feeds:", string(b))
// 	}
// }
func main() {
	// println("debug main")
	// for i := 0; i < 4; i++ {
	// 	logrus.Info(" ==== run user:", i)
	// 	server(1)
	// }
	// go getPersonnalFeed(1)
	// go celery()
	time.Sleep(10 * 6000 * time.Second)
}
