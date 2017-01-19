package g

import (
	"github.com/Sirupsen/logrus"
	"github.com/kless/osutil/user/crypt/sha512_crypt"
)

func initRoot() {
	if len(Conf.RootName) == 0 || len(Conf.RootPassword) == 0 {
		logrus.Error("There is no root user configed!")
		return
	}

	session := Session()
	db := session.DB(Conf.DBName)
	defer session.Close()
	//TODO add root user

	c := sha512_crypt.New()
	rootPassword, _ := c.Generate([]byte(Conf.RootPassword), []byte(""))
}
