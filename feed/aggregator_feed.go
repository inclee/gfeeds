package feed

import (
	"github.com/inclee/gfeeds/activity"
)

type AggregatorFeed struct {
	*BaseFeed
	aggregator Aggregator
	maxLen     int64
}

func NewAggregatorFeed(agg Aggregator, maxLen int64) *AggregatorFeed {
	return &AggregatorFeed{BaseFeed: &BaseFeed{}, aggregator: agg, maxLen: maxLen}
}

func (self *AggregatorFeed) Add(activty *activity.BaseActivty) (err error) {
	_, err = self.AddMany([]*activity.BaseActivty{activty})
	return
}

func (self *AggregatorFeed) AddMany(activties []*activity.BaseActivty) (cnt int64, err error) {
	acts, err := self.GetActivities(0, self.maxLen)
	if err != nil {
		return 0, err
	}
	newActs, oldActs := self.aggregator.Merge(acts, activties)
	_, err = self.BaseFeed.AddMany(newActs)
	if err != nil {
		return 0, err
	}
	org_len := len(acts)
	all_len := org_len + len(newActs) - len(oldActs)
	rm := int64(all_len) - self.maxLen
	rmi := 0
	var i int64
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
	return int64(len(newActs)), err
}
func (self *AggregatorFeed) Seen(actIds []int) (err error) {
	_add := make([]*activity.BaseActivty, 0, 0)
	_remove := make([]*activity.BaseActivty, 0, 0)
	acts, err := self.GetActivities(0, self.maxLen)
	if err != nil {
		return err
	}
	for _, act := range acts {
		for _, id := range actIds {
			if act.Actor == id {
				a := act.DeepCopy()
				a.Actor = 0
				_add = append(_add, a)
				_remove = append(_remove, act)
			}
		}
	}
	_, err = self.BaseFeed.AddMany(_add)
	if err == nil {
		_, err = self.BaseFeed.RemoveMany(_remove)
	}
	return
}
