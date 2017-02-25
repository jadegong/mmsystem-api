package g

import (
	"os"
	"regexp"

	"github.com/kless/osutil/user/crypt/sha512_crypt"
)

//目录或文件是否存在
func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func IsEmail(email string) bool {
	reg := regexp.MustCompile("^[a-z0-9!#$%&'*+\\/=?^_`{|}~.-]+@[a-z0-9]([a-z0-9-]*[a-z0-9])?(\\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)*$")
	return reg.Match([]byte(email))
}

func EncryptPassword(password string) string {
	c := sha512_crypt.New()
	encryptedPassword, _ := c.Generate([]byte(password), []byte(""))
	return encryptedPassword
}
