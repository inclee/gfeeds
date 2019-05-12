package feedmanager

import (
	"github.com/gocelery/gocelery"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/define"
	"github.com/inclee/gfeeds/src/feed"
)

var add_operation = func(feed feed.BaseFeed,activities []*activity.BaseActivty) {
	feed.AddMany(activities)
}

type ManagerI interface {
	GetFollowIds(user int)[]int
}

type Manager struct {
	cli *gocelery.CeleryClient
	follow_activity_limit int
}

func (m *Manager)Init()(err error){
	m.cli ,err = gocelery.NewCeleryClient(
		gocelery.NewRedisCeleryBroker("redis://"),
		gocelery.NewRedisCeleryBackend("redis://"),
		5, // number of workers
	)
	m.cli.Register("feedmanager.add_activities_operation",add_operation)
	return nil
}

func (m *Manager)getUserFeedClass()feed.Feed{
	return nil
}

func (m *Manager)getUserFollowerIds(uid int)([]int){
	panic(define.NotImplementedException)
}
func (m *Manager)getUserFeed(userid int)feed.Feed {
	feed := &feed.RedisFeed{}
	feed.Init(userid)
	return feed
}
func (m *Manager)AddActivity(uid int,activty *activity.BaseActivty)  {
	user_feed := m.getUserFeedClass()
	user_feed.Add(activty)
	followerids := m.getUserFollowerIds(uid)
	for _,fuid := range followerids{
		user_feed = m.getUserFeed(fuid)
		if _,err := m.cli.Delay("add",user_feed,[]*activity.BaseActivty{activty});err != nil{
			panic(err)
		}
	}
}