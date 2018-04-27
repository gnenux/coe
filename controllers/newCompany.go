package controllers

import (
	"github.com/astaxie/beego"
)

type NewCompanyController struct {
	beego.Controller
}

func (this *NewCompanyController) Get() {
	this.Data["Title"] = "新建企业信息"
	this.Layout = "layouts/nav.html"
	this.TplName = "company_form.html"
	this.Render()
}

func (this *NewCompanyController) Post() {

}
