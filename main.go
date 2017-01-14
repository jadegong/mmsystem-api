package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"

	"jadegong/api.mmsystem.com/g"
)

func main() {
	g.InitGlobal()

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic: %v", err)
		}

		if g.DB != nil {
			g.DB.Session.Close()
		}
		os.Remove(g.Conf.Pidfile)
	}()

	//Create pidfile
	if err := ioutil.WriteFile(g.Conf.Pidfile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644); err != nil {
		logrus.Errorf("create pidfile error: %s", err.Error())
		panic(err)
	}

	//Initialize http server
	if err := InitServer(); err != nil {
		logrus.Errorf("Init http server failed: %s", err.Error())
		panic(err)
	}
}
