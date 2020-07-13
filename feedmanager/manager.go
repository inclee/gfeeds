package feedmanager

import (
	"fmt"

	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/config"
	"github.com/inclee/gfeeds/feed"
	"github.com/inclee/gfeeds/storage"
	"github.com/inclee/gocelery"
)

func intSliceContain(slice []int, v int) bool {
	for t := range slice {
		if t == v {
			return true
		}
	}
	return false
}

type ManagerDelegate interface {
	GetFollowers(user int) []int
	GetPersonalFeed(user int) feed.Feed
	GetFeed(user int) feed.Feed
	GetInterActFeed(uid, typ int) feed.Feed
}

type Manager struct {
	cli                   *gocelery.CeleryClient
	follow_activity_limit int
	delegate              ManagerDelegate
	feeds                 *feed.RedisFeed
	badmoods              *feed.RedisFeed
}

func NewFeedManager(delegate ManagerDelegate, cfg config.ManagerConfig) *Manager {
	m := new(Manager)
	m.Init(delegate, cfg)
	m.feeds = feed.NewRedisFeed()
	m.badmoods = feed.NewRedisFeed()
	m.feeds.Init(0, fmt.Sprint("global_feed_"), storage.DefaultRedisTimelineStorage, &storage.ActiveStorage{})
	m.badmoods.Init(0, fmt.Sprint("badmoods_feed_"), storage.DefaultRedisTimelineStorage, &storage.ActiveStorage{})
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

func (m *Manager) AddInterActivity(uid, typ int, act *activity.BaseActivty) {
	if act.Target != uid {
		intactFeed := m.delegate.GetInterActFeed(act.Target, typ)
		intactFeed.Add(act)
	}
}
func (m *Manager) AddActivity(uid int, insertToSelf bool, act *activity.BaseActivty) {
	m.feeds.Add(act)
	if insertToSelf {
		user_feed := m.delegate.GetPersonalFeed(uid)
		user_feed.Add(act)

	}
	if act.Private {
		return
	}
	followerids := m.delegate.GetFollowers(uid)
	for _, fuid := range followerids {
		if len(act.Allow) > 0 {
			if intSliceContain(act.Allow, fuid) == false {
				continue
			}
		}
		if len(act.Deny) > 0 {
			if intSliceContain(act.Deny, fuid) {
				continue
			}
		}

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

func (m *Manager) LoadPersonFeeds(uid int, pgx int64, pgl int64) (acts []*activity.BaseActivty, err error) {
	user_feed := m.delegate.GetPersonalFeed(uid)
	return user_feed.GetActivities(pgx, pgl)
}

func (m *Manager) LoadInteractFeeds(uid, typ int, pgx, pgl int64) (acts []*activity.BaseActivty, newIds []int, err error) {
	_feed := m.delegate.GetInterActFeed(uid, typ)
	if actFeed, ok := _feed.(*feed.AggregatorFeed); ok {
		newIds, err = actFeed.NewIds()
		if err != nil {
			return acts, newIds, err
		}
		acts, err = _feed.GetActivities(pgx, pgl)
	}
	return acts, newIds, err
}
func (m *Manager) SeeInteractFeeds(uid, typ int, feedsId []int) {
	_feed := m.delegate.GetInterActFeed(uid, typ)
	if actFeed, ok := _feed.(*feed.AggregatorFeed); ok {
		actFeed.Seen(feedsId)
	}
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

func (m *Manager) LoadFeeds(uid int, pgx int64, pgl int64) (acts []*activity.BaseActivty, err error) {
	user_feed := m.delegate.GetFeed(uid)
	return user_feed.GetActivities(pgx, pgl)
}

func (m *Manager) LoadGlobalFeeds(pgx int64, pgl int64) (acts []*activity.BaseActivty, err error) {
	return m.feeds.GetActivities(pgx, pgl)
}
func (m *Manager) AddBadMoodFeed(act *activity.BaseActivty) (err error) {
	return m.badmoods.Add(act)
}
func (m *Manager) LoadBadMoodFeeds(pgx int64, pgl int64) (acts []*activity.BaseActivty, err error) {
	return m.badmoods.GetActivities(pgx, pgl)
}
