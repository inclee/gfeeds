package feed

type RedisFeed struct {
	*BaseFeed
}

func NewRedisFeed()*RedisFeed  {
	return &RedisFeed{
		BaseFeed:&BaseFeed{},
	}
}

func (self *RedisFeed)KeyFormat()string{
	return "feed_%ds"
}