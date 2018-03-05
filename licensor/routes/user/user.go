package user

import (
	"github.com/kataras/iris"
)

// 用户对象
type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

// 返回对象
type DataResult struct {
	Errno  int                    `json:"errno"`
	ErrMsg string                 `json:"errmsg"`
	Data   map[string]interface{} `json:"data"`
}

// 用户登录
func Login(ctx iris.Context) {
	var user User
	ctx.ReadJSON(&user)
	data := make(map[string]interface{})
	data["name"] = user.Username + "login"
	rtn := DataResult{
		Errno:  0,
		ErrMsg: "",
		Data:   data,
	}
	ctx.JSON(rtn)
}

// 用户注册
func Register(ctx iris.Context) {
	var user User
	ctx.ReadJSON(&user)
	data := make(map[string]interface{})
	data["name"] = user.Username + "register"
	rtn := DataResult{
		Errno:  0,
		ErrMsg: "",
		Data:   data,
	}
	ctx.JSON(rtn)
}

// 用户退出
func Exit(ctx iris.Context) {
	var user User
	ctx.ReadJSON(&user)
	data := make(map[string]interface{})
	data["name"] = user.Username + "exit"
	rtn := DataResult{
		Errno:  0,
		ErrMsg: "",
		Data:   data,
	}
	ctx.JSON(rtn)
}
