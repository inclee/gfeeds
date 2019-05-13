package feedmanager

import (
	"github.com/gocelery/gocelery"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/feed"
	"github.com/inclee/gfeeds/src/storage/redis"
)

type ManagerDelegate interface {
	GetFollowIds(user int)[]int
	GetPersonalFeed(user int)feed.Feed
}

type Manager struct {
	cli *gocelery.CeleryClient
	follow_activity_limit int
	delegate ManagerDelegate
}

func NewFeedManager(delegate ManagerDelegate) *Manager {
	m := new(Manager)
	m.Init(delegate)
	return m
}
func (m *Manager)Init(delegate ManagerDelegate)(err error){
	m.cli ,err = gocelery.NewCeleryClient(
		gocelery.NewRedisCeleryBroker("redis://192.168.21.231:6379"),
		gocelery.NewRedisCeleryBackend("redis://192.168.21.231:6379"),
		5, // number of workers
	)
	m.delegate = delegate
	//m.cli.Register("feedmanager.add_activities_operation",Add_operation)
	return nil
}
func (m *Manager)getUserFeed(userid int)feed.Feed {
	feed := feed.NewRedisFeed()
	feed.Init(userid,redis.NewRedisTimeLineStorage(new(redis.RedisTimeLineStorageDelegate)),&redis.ActiveStorage{})
	return feed
}
func (m *Manager)AddActivity(uid int,act*activity.BaseActivty)  {
	user_feed := m.delegate.GetPersonalFeed(uid)
	user_feed.Add(act)
	followerids := m.delegate.GetFollowIds(uid)
	for _,fuid := range followerids{
		user_feed = m.getUserFeed(fuid)
		actBytes,err := act.JsonSerialize()
		if err == nil {
			if _,err := m.cli.DelayKwargs("feedmanager.add_activities_operation", map[string]interface{}{
				"user":fuid,
				"activities": []string{string(actBytes)},
			});err != nil{
				panic(err)
			}
		}

	}
}