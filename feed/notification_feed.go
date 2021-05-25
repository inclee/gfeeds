package feed

import (
	"github.com/inclee/gfeeds/storage"
)

type NotificationFeed struct {
	*BaseFeed
}

func NewNotificationFeed(uid uint64) *NotificationFeed {
	self := &NotificationFeed{}
	self.BaseFeed = &BaseFeed{UserId: uid, TimelineStorage: storage.DefaultRedisTimelineStorage}
	return self
}

func (n *NotificationFeed) Read() {

}
