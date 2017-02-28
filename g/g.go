package g

import (
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
)

var (
	TempDir string
)

func InitGlobal() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})

	//处理panic产生的错误
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Initialization failed: %v", err)
		}
	}()
	err := initConfig()
	if err != nil {
		panic("Load config file failed.")
	}
	logrus.Infof("Service Version: %s(%s)", VERSION, Conf.Mode)

	// Set runtime procs
	runtime.GOMAXPROCS(Conf.MaxProc)

	// Initialize log
	initLogger()

	//Initialize database
	initDB()

	//Initialize storage
	initStorage()

	//Initialize cache
	InitCache()

	//TODO Initialize rpcclient

	//TODO Initialize root user
	//if errMsg := initRoot(); errMsg != "" {
	//	panic(errMsg)
	//}

	TempDir = os.TempDir()
	logrus.Infof("Temp file directory: %s", TempDir)
}
