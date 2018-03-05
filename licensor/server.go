package main

import (
	"flag"

	"github.com/kataras/iris"
	"github.com/panenming/go-im/licensor/routes"
)

var addr = flag.String("addr", ":5000", "port to listen")

func main() {
	flag.Parse()

	app := iris.New()
	app.RegisterView(iris.HTML("views", ".html").Reload(true))
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		// .Values are used to communicate between handlers, middleware.
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}

		ctx.Writef("(Unexpected) internal server error")
	})

	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})

	// 注册路由
	routes.Routes(app)

	app.Run(iris.Addr(*addr), iris.WithCharset("UTF-8"))
}
