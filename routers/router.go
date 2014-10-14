// @APIVersion 1.0.0
// @Title Pet Rest API
package routers

import (
	"pet/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/articles",
			beego.NSInclude(
				&controllers.ArticlesController{},
			),
		),

		beego.NSNamespace("/likes",
			beego.NSInclude(
				&controllers.LikesController{},
			),
		),

		beego.NSNamespace("/comments",
			beego.NSInclude(
				&controllers.PhotoCommentsController{},
			),
		),

		beego.NSNamespace("/photos",
			beego.NSInclude(
				&controllers.PhotosController{},
			),
		),

		beego.NSNamespace("/timeline",
			beego.NSInclude(
				&controllers.TimelineController{},
			),
		),

		beego.NSNamespace("/ul",
			beego.NSRouter("/follower", &controllers.UserRelationsController{}, "get:Follower"),
			beego.NSRouter("/following", &controllers.UserRelationsController{}, "get:Following"),
			beego.NSInclude(
				&controllers.UserRelationsController{},
			),
		),

		beego.NSNamespace("/users",
			beego.NSRouter("/login", &controllers.UsersController{}, "get:Login"),
			beego.NSRouter("/logout", &controllers.UsersController{}, "get:Logout"),
			beego.NSRouter("/register", &controllers.UsersController{}, "post:Register"),

			beego.NSInclude(
				&controllers.UsersController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
