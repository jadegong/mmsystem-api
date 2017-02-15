package g

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/Terry-Mao/goconf"
)

var (
	gconf          *goconf.Config
	Conf           *Config
	confFile       string
	ServiceControl string
)

type Config struct {
	//base
	Mode             string        `goconf:"base:mode"`
	Pidfile          string        `goconf:"base:pidfile"`
	MaxProc          int           `goconf:"base:maxproc"`
	HTTPBind         string        `goconf:"base:http.bind"`
	HTTPReadTimeout  time.Duration `goconf:"base:http.read.timeout:time"`
	HTTPWriteTimeout time.Duration `goconf:"base:http.write.timeout:time"`
	XSendFile        bool          `goconf:"base:http.sendfile"`
	Cache            string        `goconf:"base:cache"`
	//Storage	string	`goconf:"base.storage"`
	StoragePath      string        `goconf:"base:storage.path"`
	StorageThumbPath string        `goconf:"base:storage.thumb.path"`
	TokenDuration    time.Duration `goconf:"base:token.duration:time"`
	//db
	DBUrl     string        `goconf:"db:db.url"`
	DBName    string        `goconf:"db:db.name"`
	DBTimeout time.Duration `goconf:"db:db.timeout:time"`
	//redis
	RedisUrl      string `goconf:"redis:redis.url"`
	RedisPassword string `goconf:"redis:redis.password"`
	//rpc
	RPCProtocol        string        `goconf:"rpc:protocol"`
	RPCAddress         string        `goconf:"rpc:address"`
	RPCPoolMaxTotal    int           `goconf:"rpc:pool.maxTotal"`
	RPCPoolMaxIdle     int           `goconf:"rpc:pool.maxIdle"`
	RPCPoolMinIdle     int           `goconf:"rpc:pool.minIdle"`
	RPCPoolIdleTimeout time.Duration `goconf:"rpc:pool.idleTimeout:time"`
	//log
	LogEngine     string `goconf:"log:engine"`
	LogFilePath   string `goconf:"log:file.path"`
	LogFilePrefix string `goconf:"log:file.prefix"`
	LogMongodbUrl string `goconf:"log:mongodb.url"`
	LogMongodbDB  string `goconf:"log:mongodb.db"`
	//root user
	RootEmail    string `goconf:"root:email"`
	RootPassword string `goconf:"root:password"`
}

func newConfig() *Config {
	return &Config{
		Mode:               "dev",
		Pidfile:            "/var/run/mmsystem.pid",
		MaxProc:            runtime.NumCPU(),
		HTTPBind:           "0.0.0.0:8088",
		HTTPReadTimeout:    5 * time.Second,
		HTTPWriteTimeout:   5 * time.Second,
		XSendFile:          true,
		Cache:              "redis",
		StoragePath:        "/data/mmsystem/files",
		StorageThumbPath:   "/data/mmsystem/thumbs",
		TokenDuration:      2 * time.Hour,
		DBUrl:              "mongodb://127.0.0.1:27017",
		DBName:             "mmsystem",
		DBTimeout:          3 * time.Second,
		RedisUrl:           "127.0.0.1:6379",
		RedisPassword:      "",
		RPCProtocol:        "binary",
		RPCAddress:         "127.0.0.1:8081",
		RPCPoolMaxTotal:    100,
		RPCPoolMaxIdle:     10,
		RPCPoolMinIdle:     5,
		RPCPoolIdleTimeout: 2 * time.Hour,
		LogEngine:          "mongodb",
		LogFilePath:        "/data/mmsystem/log",
		LogFilePrefix:      "mmsystem",
		LogMongodbUrl:      "mongodb://127.0.0.1:27017",
		LogMongodbDB:       "mmsystem_log",
		RootEmail:          "root@mmsystem.com",
		RootPassword:       "root@mmsystem",
	}
}

func init() {
	flag.StringVar(&confFile, "c", "./config.conf", " set config file path")
	flag.Parse()
}

func initConfig() (err error) {
	Conf = newConfig()
	gconf = goconf.New()
	if err = gconf.Parse(confFile); err != nil {
		return err
	}
	if err := gconf.Unmarshal(Conf); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func ReloadConfig() (*Config, error) {
	conf := newConfig()
	ngconf, err := gconf.Reload()
	if err != nil {
		return nil, err
	}
	if err = ngconf.Unmarshal(conf); err != nil {
		return nil, err
	}
	gconf = ngconf
	return conf, nil
}
