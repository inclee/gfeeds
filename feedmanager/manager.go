package feedmanager

import (
	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/config"
	"github.com/inclee/gfeeds/feed"
	"github.com/inclee/gocelery"
)

type ManagerDelegate interface {
	GetFollowers(user int) []int
	GetPersonalFeed(user int) feed.Feed
	GetFeed(user int) feed.Feed
}

type Manager struct {
	cli                   *gocelery.CeleryClient
	follow_activity_limit int
	delegate              ManagerDelegate
	feeds                 feed.Feed
}

func NewFeedManager(delegate ManagerDelegate, cfg config.ManagerConfig) *Manager {
	m := new(Manager)
	m.Init(delegate, cfg)
	m.feeds = feed.NewRedisFeed()
	config.Config = cfg
	return m
}
func (m *Manager) Init(delegate ManagerDelegate, cfg config.ManagerConfig) (err error) {
	m.cli, err = gocelery.NewCeleryClient(
		gocelery.NewRedisCeleryBroker(cfg.CeleryBroker),
		gocelery.NewRedisCeleryBackend(cfg.CeleryBackend),
		cfg.CeleryWorkNum, // number of workers
	)
	m.delegate = delegate
	return nil
}

func (m *Manager) AddActivity(uid int, act *activity.BaseActivty) {
	m.feeds.Add(act)
	user_feed := m.delegate.GetPersonalFeed(uid)
	user_feed.Add(act)
	followerids := m.delegate.GetFollowers(uid)
	for _, fuid := range followerids {
		actBytes, err := act.Serialize()
		if err == nil {
			if _, err := m.cli.DelayKwargs("feedmanager.add_activities_operation", map[string]interface{}{
				"user":       fuid,
				"activities": []string{string(actBytes)},
			}); err != nil {
				panic(err)
			}
		}
	}

}

func (m *Manager) LoadPersonFeeds(uid int, pgx int, pgl int) []*activity.BaseActivty {
	user_feed := m.delegate.GetPersonalFeed(uid)
	return user_feed.GetActivities(pgx, pgl)
}

func (m *Manager) InsertFeedActivities(uid int, acts []*activity.BaseActivty) {
	for _, act := range acts {
		actBytes, err := act.Serialize()
		if err == nil {
			if _, err := m.cli.DelayKwargs("feedmanager.remove_activities_operation", map[string]interface{}{
				"user":       uid,
				"activities": []string{string(actBytes)},
			}); err != nil {
				panic(err)
			}
		}
	}
}

func (m *Manager) RemoveFeedActivities(uid int, acts []*activity.BaseActivty) {
	for _, act := range acts {
		actBytes, err := act.Serialize()
		if err == nil {
			if _, err := m.cli.DelayKwargs("feedmanager.remove_activities_operation", map[string]interface{}{
				"user":       uid,
				"activities": []string{string(actBytes)},
			}); err != nil {
				panic(err)
			}
		}
	}
}

func (m *Manager) LoadFeeds(uid int, pgx int, pgl int) []*activity.BaseActivty {
	return m.feeds.GetActivities(pgx, pgl)
}
