package feed

import (
	"fmt"
	"strconv"

	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/storage"
)

type AggregatorFeed struct {
	*NotificationFeed
	aggregator Aggregator
	id         string
	maxLen     int
}

func NewAggregatorFeed(uid uint64, id string, agg Aggregator, maxLen int) *AggregatorFeed {
	return &AggregatorFeed{NotificationFeed: NewNotificationFeed(uid), aggregator: agg, id: id, maxLen: maxLen}
}

func (self *AggregatorFeed) Add(activty *activity.BaseActivty) (err error) {
	if activty.Actor == self.UserId {
		return nil
	}
	_, err = self.AddMany([]*activity.BaseActivty{activty})
	return
}

func (self *AggregatorFeed) AddMany(activties []*activity.BaseActivty) (cnt int64, err error) {
	acts, err := self.BaseFeed.GetActivities(0, self.maxLen)
	if err != nil {
		return 0, err
	}
	newActs, oldActs := self.aggregator.Merge(acts, activties)
	if len(newActs) > 0 {
		_, err = self.BaseFeed.AddMany(newActs)
		if err != nil {
			return 0, err
		}
	}

	org_len := len(acts)
	all_len := org_len + len(newActs) - len(oldActs)
	rm := all_len - self.maxLen
	rmi := 0
	var i int
	for ; i < rm; rmi++ {
		oact := acts[org_len-rmi-1]
		find := false
		for _, o := range oldActs {
			if o == oact {
				find = true
				break
			}
		}
		if !find {
			oldActs = append(oldActs, oact)
			i++
		}
	}
	if len(oldActs) > 0 {
		_, err = self.BaseFeed.RemoveMany(oldActs)
	}
	newId := make([]string, 0, 0)
	for _, na := range newActs {
		newId = append(newId, strconv.Itoa(na.VerbObj.Id))
	}
	err = storage.RedisClent.SAdd(fmt.Sprintf("interact_%d_%s", self.UserId, self.id), newId).Err()
	return int64(len(newActs)), err
}

// 	if actFeed, ok := _feed.(*feed.AggregatorFeed); ok {
// 		newIds, err = actFeed.NewIds()
// 		if err != nil {
// 			return acts, newIds, err
// 		}
// 		acts, err = _feed.GetActivities(pgx, pgl)
// 	}
// 	return acts, newIds, err
func (self *AggregatorFeed) GetActivities(pgx int, pgl int) (acts []*activity.BaseActivty, newIds []int, err error) {
	acts, err = self.BaseFeed.GetActivities(pgx, pgl)
	if err != nil {
		return
	}
	newIds, err = self.NewIds()
	return
}

func (self *AggregatorFeed) NewIds() ([]int, error) {
	ret := make([]int, 0, 0)
	ids, err := storage.RedisClent.SMembers(fmt.Sprintf("interact_%d_%s", self.UserId, self.id)).Result()
	if err != nil {
		return []int{}, err
	}
	for _, id := range ids {
		aid, _ := strconv.Atoi(id)
		ret = append(ret, aid)
	}
	return ret, err
}

func (self *AggregatorFeed) Seen(actIds []int) (err error) {
	// err = storage.RedisClent.SRem(fmt.Sprintf("interact_%d_%d", self.UserId, self.typ), actIds).Err()
	// _add := make([]*activity.BaseActivty, 0, 0)
	// _remove := make([]*activity.BaseActivty, 0, 0)
	// acts, err := self.GetActivities(0, self.maxLen)
	// if err != nil {
	// 	return err
	// }
	// for _, act := range acts {
	// 	for _, id := range actIds {
	// 		if act.Actor == id {
	// 			a := act.DeepCopy()
	// 			a.Actor = 0
	// 			_add = append(_add, a)
	// 			_remove = append(_remove, act)
	// 		}
	// 	}
	// }
	// _, err = self.BaseFeed.AddMany(_add)
	// if err == nil {
	// 	_, err = self.BaseFeed.RemoveMany(_remove)
	// }
	return err
}
