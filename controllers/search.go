package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gnenux/coe/models"
)

type SearchController struct {
	beego.Controller
}

type BasicCompany struct {
	ID   int
	Name string
}

func (this *SearchController) Get() {
	this.Data["Title"] = "查询结果"
	this.Data["URL"] = URL
	this.Layout = "layouts/nav.html"

	keys := this.GetStrings("key")
	if len(keys) == 0 {
		this.Data["Error"] = "not found key!"
		this.TplName = "error.html"
	} else {
		companies, err := models.GetCompaniesByKeys(keys)
		if err != nil {
			this.Data["Error"] = err.Error()
			this.TplName = "error.html"
		} else {
			if this.Ctx.Input.AcceptsJSON() {
				bcs := []BasicCompany{}
				for _, v := range *companies {
					bcs = append(bcs, BasicCompany{v.ID, v.Name})
				}
				this.Data["json"] = &bcs
				this.ServeJSON()
				return
			}
			this.Data["Companies"] = companies
			this.TplName = "companies.html"
		}

	}

	this.Render()
}
