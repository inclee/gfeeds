package config

type ManagerConfig struct {
	CeleryBroker    string
	CeleryBackend   string
	CeleryWorkNum   int
	TimelineStorage []string
}

var Config ManagerConfig
