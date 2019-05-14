package config

type CeleryConf struct {
	Broker string
	Backend string
	WorkNum int
}
type TimelineStorageConf struct {
	Redis []string
}
type Config struct {
	Celery	CeleryConf
	TimelineStorage TimelineStorageConf
}

