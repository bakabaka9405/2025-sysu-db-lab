package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrForbidden           = newError(403, "Forbidden")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// user biz errors
	ErrUsernameExists     = newError(1002, "用户名已存在")
	ErrPhoneExists        = newError(1003, "手机号已被注册")
	ErrEmailExists        = newError(1004, "邮箱已被注册")
	ErrInvalidCredentials = newError(1005, "用户名或密码错误")
)
