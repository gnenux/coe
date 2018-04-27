package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "AddBadManagement",
			Router: `/:id/badmanagements`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "RemoveBadManagement",
			Router: `/:id/badmanagements/:bid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "GetBadManagement",
			Router: `/:id/badmanagements/:bid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "AddBlackList",
			Router: `/:id/blacklist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "RemoveBlackList",
			Router: `/:id/blacklist/:bid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "GetBlackList",
			Router: `/:id/blacklist/:bid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "GetPenalty",
			Router: `/:id/penalties/:pid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "AddPenalty",
			Router: `/:id/penalty`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"] = append(beego.GlobalControllerRouter["github.com/gnenux/coe/controllers:CompaniesController"],
		beego.ControllerComments{
			Method: "RemovePenalty",
			Router: `/:id/penalty/:pid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

}
