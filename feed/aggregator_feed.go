package feed

import (
	"github.com/inclee/gfeeds/activity"
)

type AggregatorFeed struct {
	*BaseFeed
	aggregator Aggregator
	maxLen     int
}

func NewAggregatorFeed(agg Aggregator, maxLen int) *AggregatorFeed {
	return &AggregatorFeed{BaseFeed: &BaseFeed{}, aggregator: agg, maxLen: maxLen}
}

func (self *AggregatorFeed) Add(activty *activity.BaseActivty) {
	self.AddMany([]*activity.BaseActivty{activty})
}

func (self *AggregatorFeed) AddMany(activties []*activity.BaseActivty) int {
	acts := self.GetActivities(0, self.maxLen)
	newActs, oldActs := self.aggregator.Merge(acts, activties)
	self.BaseFeed.AddMany(newActs)
	org_len := len(acts)
	all_len := org_len + len(newActs) - len(oldActs)
	rm := all_len - self.maxLen
	rmi := 0
	for i := 0; i < rm; rmi++ {
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
			i += 1
		}
	}
	if len(oldActs) > 0 {
		self.BaseFeed.RemoveMany(oldActs)
	}
	return len(newActs)
}
func (self *AggregatorFeed) Seen(actIds []int) {
	_add := make([]*activity.BaseActivty, 0, 0)
	_remove := make([]*activity.BaseActivty, 0, 0)
	acts := self.GetActivities(0, self.maxLen)
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
	self.BaseFeed.AddMany(_add)
	self.BaseFeed.RemoveMany(_remove)
}
