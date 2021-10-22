package main

import (
	"flag"
	"strings"
	"time"
)

type IConfig interface {
	Mode() string
	TimeFrameWidth() int64
	UpdatePeriod() int64
}

type Config struct {
	mode           string
	timeFrameWidth int64
	updatePeriod   int64
}

func NewConfig() IConfig {
	c := &Config{}
	c.init()
	return c
}

func (c *Config) init() {
	c.updatePeriod = *flag.Int64("upd", int64(1*time.Second), "update period (nanosec) default: 1s (1000000000 ns)")
	c.mode = strings.ToUpper(*flag.String("mode", "ALL", "set data source server : JSON or SSE or ALL (default)"))
	c.timeFrameWidth = *flag.Int64("tf", int64(10*time.Second), "width TimeFrame (ns) int, default : 10 s (10000000000 ns) ")
	flag.Parse()
}

func (c *Config) TimeFrameWidth() int64 {
	return c.timeFrameWidth
}

func (c *Config) Mode() string {
	return c.mode
}

func (c *Config) UpdatePeriod() int64 {
	return c.updatePeriod
}
