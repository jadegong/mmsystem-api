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
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2010-01-02 12:34:56.00"})
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

	//TODO Initialize rpcclient

	//Initialize root user
	initRoot()

	TempDir = os.TempDir()
	logrus.Infof("Temp file directory: %s", TempDir)
}
