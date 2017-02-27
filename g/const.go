package g

const (
	DOMAIN     = "https://mmsystem.com"
	DOMAIN_API = "https://api.mmsystem.com"
	VERSION    = "0.0.1"
	AUTH_REALM = "mmsystem.com"

	AUTH_TOKEN_NS = "auth:" //认证token存取前缀

	CACHE_REDIS  = "redis"
	CACHE_MEMORY = "memory"

	RUN_MODE_DEV     = "dev"
	RUN_MODE_PRODUCT = "prod"

	LOG_ENGINE_FILE    = "file"
	LOG_ENGINE_MONGODB = "mongodb"

	THUMB_FOLDER_NUM = 1024
	THUMB_FORMAT     = ".png"

	LOCAL_TIME_ZONE_OFFSET = 8 * 60 * 60 //Beijing(UTC+8:00)

	DEFAULT_FILE_MODE = 0755

	USER = "user"
)

//todo log type
const (
	LOG_TYPE_DEFAULT = iota
	LOG_TYPE_REGISTER
	LOG_TYPE_ACTIVATION
	LOG_TYPE_LOGIN
)
