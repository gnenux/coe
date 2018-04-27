package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Title"] = "企业信用查询系统"
	this.Data["URL"] = URL
	this.Layout = "layouts/nav.html"
	this.TplName = "index.html"
	this.Render()
}
