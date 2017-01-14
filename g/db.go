package g

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/lunny/log"
	"gopkg.in/mgo.v2"
)

var (
	DB           *mgo.Database
	MongoSession *mgo.Session
)

type DBLogger struct {
	logger *logrus.Logger
}

func (l *DBLogger) Output(calldepth int, s string) error {
	l.logger.Print(s)
	return nil
}

func NewDBLogger() *DBLogger {
	return &DBLogger{
		logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: "2010-01-01 11:22:11",
			},
			Hooks: make(logrus.LevelHooks),
			Level: logrus.InfoLevel,
		},
	}
}

func initDB() {
	var err error
	if Conf.Mode == RUN_MODE_DEV {
		mgo.SetDebug(true)
		logger := NewDBLogger()
		mgo.SetLogger(logger)
	}

	MongoSession, err = mgo.DialWithTimeout(Conf.DBUrl, Conf.DBTimeout)
	if err != nil {
		panic(err)
	}
	MongoSession.SetMode(mgo.Monotonic, true)
	DB = MongoSession.DB(Conf.DBName)
}

func Session() *mgo.Session {
	if MongoSession == nil {
		initDB()
	}
	return MongoSession.Clone()
}

func M(collection string, f func(*mgo.Collection)) {
	session := Session()
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			log.Errorf("MongoDB error: %v", err)
		}
	}()
	c := session.DB(Conf.DBName).C(collection)
	f(c)
}
