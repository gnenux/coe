// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/gnenux/coe/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	ns := beego.NewNamespace("/companies", beego.NSInclude(&controllers.CompaniesController{}))
	beego.AddNamespace(ns)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/search", &controllers.SearchController{})
	beego.Router("/newcompany", &controllers.NewCompanyController{})

	beego.InsertFilter("/static/*", beego.BeforeStatic, func(ctx *context.Context) {
		ctx.Output.Header("Cache-Control", "max-age=3600")
	})

}
