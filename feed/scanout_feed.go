package feed

import (
	"fmt"

	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/storage"
	"github.com/inclee/gfeeds/util"
	"github.com/inclee/gocelery"
)

type ScanOutFeed struct {
	user   *BaseFeed
	global *BaseFeed
	cli    *gocelery.CeleryClient
}

func NewScanOutFeed(uid uint64, id string, cli *gocelery.CeleryClient) *ScanOutFeed {
	self := &ScanOutFeed{}
	self.user = &BaseFeed{UserId: uid}
	self.global = &BaseFeed{}
	self.cli = cli
	self.user.Init(uid, fmt.Sprintf("%d_%s", uid, id), storage.DefaultRedisTimelineStorage, &storage.ActiveStorage{})
	self.global.Init(uid, fmt.Sprintf("global_%s", id), storage.DefaultRedisTimelineStorage, &storage.ActiveStorage{})
	return self
}

func (f *ScanOutFeed) Add(act *activity.BaseActivty, outUids []uint64, targetFeedId string) error {
	f.user.Add(act)
	if act.Private {
		return nil
	}
	f.global.Add(act)
	for _, fuid := range outUids {
		if len(act.Allow) > 0 {
			if util.UInt64SliceContain(act.Allow, fuid) == false {
				continue
			}
		}
		if len(act.Deny) > 0 {
			if util.UInt64SliceContain(act.Deny, fuid) {
				continue
			}
		}
		actBytes, err := act.Serialize()
		if err == nil {
			if _, err := f.cli.DelayKwargs("feedmanager.add_activities_operation", map[string]interface{}{
				"user":           fuid,
				"activities":     []string{string(actBytes)},
				"target_feed_id": targetFeedId,
			}); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func (f *ScanOutFeed) AddMany(acts []*activity.BaseActivty, outUids []uint64, targetFeedId string) error {
	f.user.AddMany(acts)
	for _, act := range acts {
		if act.Private {
			return nil
		}
		f.global.Add(act)
		for _, fuid := range outUids {
			if len(act.Allow) > 0 {
				if util.UInt64SliceContain(act.Allow, fuid) == false {
					continue
				}
			}
			if len(act.Deny) > 0 {
				if util.UInt64SliceContain(act.Deny, fuid) {
					continue
				}
			}

			actBytes, err := act.Serialize()
			if err == nil {
				if _, err := f.cli.DelayKwargs("feedmanager.add_activities_operation", map[string]interface{}{
					"user":           fuid,
					"activities":     []string{string(actBytes)},
					"target_feed_id": targetFeedId,
				}); err != nil {
					panic(err)
				}
			}
		}
	}
	return nil
}
