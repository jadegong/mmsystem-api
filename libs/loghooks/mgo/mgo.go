package mgo

import (
	"fmt"
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	LOG_DAILY = iota
	LOG_MONTHLY
	LOG_YEARLY
)

// 级别配置，按周期自动分表monthly、yearly、daily {log.InfoLevel: {"period": "monthly", "prefix"：“”}}
type LevelMap map[logrus.Level]map[string]interface{}

type mgoHook struct {
	db     *mgo.Database
	cMap   LevelMap
	levels []logrus.Level
}

type M bson.M

func NewHook(mgoUrl, db string, levelMap LevelMap) *mgoHook {
	session, err := mgo.Dial(mgoUrl)
	if err != nil {
		panic(err)
	}

	hook := &mgoHook{db: session.DB(db), cMap: levelMap}
	for level, _ := range levelMap {
		hook.levels = append(hook.levels, level)
	}
	return hook
}

func (h *mgoHook) Fire(entry *logrus.Entry) error {
	conf, ok := h.cMap[entry.Level]
	if !ok {
		err := fmt.Errorf("no collection provided for loglevel: %d", entry.Level)
		log.Println(err.Error())
		return err
	}
	collection := entry.Level.String() + "_"
	prefix, _ := conf["prefix"].(string)
	if prefix != "" {
		collection = prefix + collection
	}
	now := time.Now()
	switch conf["period"].(int) {
	case LOG_DAILY:
		collection += now.Format("20060102")
	case LOG_MONTHLY:
		collection += now.Format("200601")
	case LOG_YEARLY:
		collection += now.Format("2006")
	}
	entry.Data["Level"] = entry.Level.String()
	entry.Data["Time"] = entry.Time
	entry.Data["Message"] = entry.Message
	err := h.db.C(collection).Insert(M(entry.Data))
	if err != nil {
		return fmt.Errorf("Failed to send log entry to mongodb: %s", err)
	}

	return nil
}

func (h *mgoHook) Levels() []logrus.Level {
	return h.levels
}
