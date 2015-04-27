// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"blog-on-go/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/comments",
			beego.NSInclude(
				&controllers.CommentsController{},
			),
		),

		beego.NSNamespace("/links",
			beego.NSInclude(
				&controllers.LinksController{},
			),
		),

		beego.NSNamespace("/options",
			beego.NSInclude(
				&controllers.OptionsController{},
			),
		),

		beego.NSNamespace("/posts",
			beego.NSInclude(
				&controllers.PostsController{},
			),
		),

		beego.NSNamespace("/terms",
			beego.NSInclude(
				&controllers.TermsController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
