package routes

import (
	"github.com/kataras/iris"
	"github.com/panenming/go-im/licensor/routes/user"
)

// 注册路由
func Routes(app *iris.Application) {
	// login
	app.Post("/user/login", user.Login)
	// 注册
	app.Post("/user/register", user.Register)
	// 退出
	app.Post("/user/exit", user.Exit)
}
