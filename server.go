package main

import "jadegong/api.mmsystem.com/g"

func InitServer() error {
	e := initRouter()
	e.ReadTimeout = g.Conf.HTTPReadTimeout
	e.WriteTimeout = g.Conf.HTTPWriteTimeout
	if err := e.Start(g.Conf.HTTPBind); err != nil {
		return err
	}
	return nil
}
