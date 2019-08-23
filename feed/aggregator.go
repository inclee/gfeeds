package feed

import (
	"encoding/json"
	"strconv"
	"time"

	set "github.com/deckarep/golang-set"
	"github.com/inclee/gfeeds/activity"
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
			if _cur.Verb == _new.Verb && _cur.Object.Id == _new.Object.Id && _cur.Object.Type == _new.Object.Type && _cur.Actor > 0 { //find
				actSet := set.NewSet()
				json.Unmarshal([]byte(_cur.Extra), &actSet)
				actor := strconv.Itoa(_new.Actor)
				if actSet.Contains(actor) == false {
					a := _cur.DeepCopy()
					actSet.Add(actor)
					if data, err := json.Marshal(actSet); err == nil {
						a.Extra = string(data)
						_add = append(_add, a)
						_remove = append(_remove, _cur)
					}
				}
				find = true
			}
		}
		if find == false {
			actset := set.NewSet(strconv.Itoa(_new.Actor))
			if extra, err := json.Marshal(actset); err == nil {
				_newInact := activity.NewActivity(int(time.Now().Unix()), _new.Verb, _new.Object, _new.Target, _new.Time, false, []int{}, []int{}, string(extra))
				_add = append(_add, _newInact)
			}
		}
	}
	return _add, _remove
}
