package feedmanager

import (
	"github.com/gocelery/gocelery"
	"github.com/inclee/gfeeds/src/activity"
	"github.com/inclee/gfeeds/src/feed"
)

type ManagerDelegate interface {
	GetFollowIds(user int)[]int
	GetPersonalFeed(user int)feed.Feed
	GetFeed(user int)feed.Feed
}

type Manager struct {
	cli *gocelery.CeleryClient
	follow_activity_limit int
	delegate ManagerDelegate
}

type ManagerConfig struct {
	CeleryBroker string
	CeleryBackend string
	CeleryWorkNum int
	TimelineStorage []string
}
var Config ManagerConfig
func NewFeedManager(delegate ManagerDelegate,config ManagerConfig) *Manager {
	m := new(Manager)
	m.Init(delegate,config)
	Config = config
	return m
}
func (m *Manager)Init(delegate ManagerDelegate,config ManagerConfig)(err error){
	m.cli ,err = gocelery.NewCeleryClient(
		gocelery.NewRedisCeleryBroker(config.CeleryBroker),
		gocelery.NewRedisCeleryBackend(config.CeleryBackend),
		config.CeleryWorkNum, // number of workers
	)
	m.delegate = delegate
	return nil
}

func (m *Manager)AddActivity(uid int,act*activity.BaseActivty)  {
	user_feed := m.delegate.GetPersonalFeed(uid)
	user_feed.Add(act)
	followerids := m.delegate.GetFollowIds(uid)
	followerids = append(followerids, uid)
	for _,fuid := range followerids{
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

func (m *Manager)LoadPersonFeeds(uid int,pgx int,pgl int)[]*activity.BaseActivty{
	user_feed := m.delegate.GetPersonalFeed(uid)
	return user_feed.GetActivities(pgx,pgl)
}

func (m *Manager)InsertFeedActivities(uid int,acts []*activity.BaseActivty)  {
	for _,act := range acts{
		actBytes,err := act.JsonSerialize()
		if err == nil {
			if _,err := m.cli.DelayKwargs("feedmanager.add_activities_operation", map[string]interface{}{
				"user":uid,
				"activities": []string{string(actBytes)},
			});err != nil{
				panic(err)
			}
		}
	}
}

func (m *Manager)LoadFeeds(uid int,pgx int,pgl int)[]*activity.BaseActivty{
	user_feed := m.delegate.GetFeed(uid)
	return user_feed.GetActivities(pgx,pgl)
}