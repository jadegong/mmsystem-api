package g

import "net/http"

const (
	SUCCESS = 200

	//Common
	ERR_DB_FAILED           = 10001
	ERR_EMAIL_NOT_AVAILABLE = 10002
	ERR_UPLOAD_FAILED       = 10003
	ERR_DATA_INVALID        = 10004
	ERR_NO_PERMISSIONS      = 10005
	ERR_UNAUTHORIZED        = 10006
	ERR_RESOURCE_NOT_FOUND  = 10007
	ERR_RPC_FAILED          = 10008
	ERR_THUMB_FAILED        = 10009
	ERR_REDIS_NOT_AVAILABLE = 10010

	//Register
	ERR_REG_EMAIL_EXISTS     = 20001
	ERR_REG_EMAIL_INVALID    = 20002
	ERR_REG_PASSWORD_INVALID = 20003

	//login
	ERR_USER_NOT_EXISTS   = 20101
	ERR_PASSWORD_INVALID  = 20102
	ERR_USER_NOT_VERIFIED = 20103
	ERR_LOGIN_FAILED      = 20104
	ERR_LOGOUT_FAILED     = 20105

	//reset
	ERR_FORGET_FAILED = 20301
	ERR_RESET_INVALID = 20302
	ERR_RESET_FAILED  = 20303
)

var errorText = map[int]string{
	//Common
	ERR_DB_FAILED:           "数据库错误",
	ERR_EMAIL_NOT_AVAILABLE: "邮箱地址无效",
	ERR_UPLOAD_FAILED:       "上传文件失败",
	ERR_DATA_INVALID:        "数据无效",
	ERR_NO_PERMISSIONS:      "没有操作权限",
	ERR_UNAUTHORIZED:        "未授权请求",
	ERR_RESOURCE_NOT_FOUND:  "资源未找到",
	ERR_RPC_FAILED:          "内部服务调用出错",
	ERR_THUMB_FAILED:        "文件生成缩略图失败",
	ERR_REDIS_NOT_AVAILABLE: "系统故障:缓存服务不可用",

	//Register
	ERR_REG_EMAIL_EXISTS:     "邮箱已经存在",
	ERR_REG_EMAIL_INVALID:    "邮箱地址无效",
	ERR_REG_PASSWORD_INVALID: "密码无效",

	//login
	ERR_USER_NOT_EXISTS:   "用户不存在",
	ERR_PASSWORD_INVALID:  "用户名或密码错误",
	ERR_USER_NOT_VERIFIED: "用户未审核",
	ERR_LOGIN_FAILED:      "登陆失败，认证服务不可用",
	ERR_LOGOUT_FAILED:     "注销失败",

	//reset
	ERR_FORGET_FAILED: "申请找回密码失败",
	ERR_RESET_INVALID: "找回密码请求已过期",
	ERR_RESET_FAILED:  "重置密码失败",
}

func GetErrMsg(errCode int) string {
	msg := errorText[errCode]
	if msg == "" {
		msg = "未知错误"
	}
	return msg
}

func GetErrHttpStatus(errCode int) int {
	switch errCode {
	case ERR_RPC_FAILED, ERR_DB_FAILED, ERR_REDIS_NOT_AVAILABLE:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
