package feedmanager

import "github.com/inclee/gfeeds/src/activity"

type ManagerI interface {
	GetFollowIds(user int)[]int
}

type Manager struct {

}

func (m *Manager)AddActivity(uid int,activty activity.BaseActivty)  {

}