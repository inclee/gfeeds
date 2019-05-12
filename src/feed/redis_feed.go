package feed

type RedisFeed struct {
	*BaseFeed
}

func (self *RedisFeed)KeyFormat()string{
	return "feed_%ds"
}