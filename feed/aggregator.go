package feed

import (
	"encoding/json"
	"time"

	"github.com/inclee/gfeeds/activity"
	"github.com/inclee/gfeeds/util"
)

type Aggregator interface {
	Merge([]*activity.BaseActivty, []*activity.BaseActivty) ([]*activity.BaseActivty, []*activity.BaseActivty)
}

type InteractAggregator struct {
}

func (agg *InteractAggregator) Merge(cur []*activity.BaseActivty, new []*activity.BaseActivty) ([]*activity.BaseActivty, []*activity.BaseActivty) {
	_add := make([]*activity.BaseActivty, 0, 0)
	_remove := make([]*activity.BaseActivty, 0, 0)
	for _, _new := range new {
		find := false
		for _, _cur := range cur {
			if _cur.Verb == _new.Verb && _cur.VerbObj.Equal(_new.VerbObj) && _cur.TargetObj.Equal(_new.TargetObj) && _cur.Actor > 0 { //find
				actList := make([]int, 0, 0)
				json.Unmarshal([]byte(_cur.Extra), &actList)
				actor := _new.Actor
				if util.IntSliceContain(actList, actor) == false {
					a := _cur.DeepCopy()
					actList = append(actList, actor)
					if data, err := json.Marshal(actList); err == nil {
						a.Extra = string(data)
						_add = append(_add, a)
						_remove = append(_remove, _cur)
					}
				} else {
					a := _cur.DeepCopy()
					util.IntSliceMoveTo(&actList, actor, 0)
					if data, err := json.Marshal(actList); err == nil {
						a.Extra = string(data)
						_add = append(_add, a)
						_remove = append(_remove, _cur)
					}
				}
				find = true
			}
		}
		if find == false {
			if extra, err := json.Marshal([]int{_new.Actor}); err == nil {
				_newInact := activity.NewActivity(int(time.Now().Unix()), _new.Verb, _new.VerbObj, _new.Target, _new.TargetObj, _new.Time, false, []int{}, []int{}, string(extra), _new.Context)
				_add = append(_add, _newInact)
			}
		}
	}
	return _add, _remove
}
