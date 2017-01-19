package g

import (
	"os"

	"github.com/Sirupsen/logrus"
	"jadegong/api.mmsystem.com/libs/loghooks/file"
	"jadegong/api.mmsystem.com/libs/loghooks/mgo"
)

func initLogger() {
	if Conf.Mode == RUN_MODE_DEV {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		if Conf.LogEngine == LOG_ENGINE_FILE { //file engine
			logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
			if ok, _ := IsFileExists(Conf.LogFilePath); !ok {
				os.MkdirAll(Conf.LogFilePath, 0755)
			}
			logrus.AddHook(file.NewHook(Conf.LogFilePath + Conf.LogFilePrefix + ".log"))
		} else if Conf.LogEngine == LOG_ENGINE_MONGODB { //mongodb engine
			logrus.AddHook(mgo.NewHook(Conf.LogMongodbUrl, Conf.LogMongodbDB, mgo.LevelMap{
				logrus.InfoLevel:  map[string]interface{}{"period": mgo.LOG_YEARLY},
				logrus.WarnLevel:  map[string]interface{}{"period": mgo.LOG_MONTHLY},
				logrus.ErrorLevel: map[string]interface{}{"period": mgo.LOG_MONTHLY},
			}))
		}
	}
}
