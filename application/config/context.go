package config

import (
	"sync"

	"github.com/cocktailcloud/acloud-alarm-collector/application/database"
	"github.com/cocktailcloud/acloud-alarm-collector/application/structure"
	"github.com/cocktailcloud/acloud-alarm-collector/application/utils"
)

type Context struct {
	Config *Config
	//AlarmConfig *AlarmConfig
	DBAdapter  database.DBAdapter
	MainQueue  utils.AlarmQueue
	RetryQueue utils.AlarmQueue
	AliveAlarm map[string](map[string]structure.AliveAlarm)
	Mutex      *sync.Mutex
}

func (c *Context) Init() {
	c.Mutex = &sync.Mutex{}
	c.MainQueue = utils.NewQueue()
	c.RetryQueue = utils.NewQueue()
	c.AliveAlarm = make(map[string](map[string]structure.AliveAlarm))
}

func (c *Context) GetAliveAlarm() map[string](map[string]structure.AliveAlarm) {
	c.Mutex.Lock()
	alarms := c.AliveAlarm
	c.Mutex.Unlock()
	return alarms
}

func (c *Context) SetAliveAlarm(alarms map[string](map[string]structure.AliveAlarm)) {
	c.Mutex.Lock()
	c.AliveAlarm = alarms
	c.Mutex.Unlock()
}
