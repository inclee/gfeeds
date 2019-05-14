package main

import (
	"github.com/inclee/gfeeds/src/config"
	"github.com/inclee/gfeeds/src/storage/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)



func parseConfig() {
	GConfig = new(config.Config)
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, GConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println("conf", GConfig)
}
func prepareRedisTimeline()  {
	for _,r := range  GConfig.TimelineStorage.Redis{
		cache := new(redis.RedisTimeLineCache)
		cache.Init(r)
		RedisTimeLineCaches = append(RedisTimeLineCaches, cache)
	}
}

func Prepare()  {
	parseConfig()
	prepareRedisTimeline()
}