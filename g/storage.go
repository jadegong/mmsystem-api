package g

import "os"

func initStorage() {
	if ok, _ := IsFileExists(Conf.StoragePath); !ok {
		os.MkdirAll(Conf.StoragePath, 0755)
	}
	if ok, _ := IsFileExists(Conf.StorageThumbPath); !ok {
		os.MkdirAll(Conf.StorageThumbPath, 0755)
	}
}
