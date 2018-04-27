package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/gnenux/coe/models"
)

type CompaniesController struct {
	beego.Controller
}

// @Title CreateEnterprise
// @Description 创建企业信息
// @Param body body models.Company true "body for enterprise content"
// @Success 201 {int} models.Company.Id
// @Failure 403 body is empty
// @router / [post]
func (this *CompaniesController) Create() {
	c := models.Company{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &c); err != nil {
		this.CustomAbort(403, err.Error())
		return

	}

	id, err := models.CreateCompany(&c)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}
	fmt.Println(id)
	this.Ctx.WriteString(id)
	this.Ctx.ResponseWriter.WriteHeader(201)

}

// @Title GetEnterprise
// @Description 查询企业信息
// @Param id path string true "企业id"
// @Success 200 {object} 企业json
// @Failure 403 body is empty
// @router /:id [get]
func (this *CompaniesController) Get() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if id <= 0 {
		this.Data["json"] = map[string]string{"error": "id is less zero or not found id"}
		this.ServeJSON()
		return
	}
	c, err := models.GetCompanyByID(id)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
	}

	bm, _ := models.GetBadManagementsByCompanyID(id)
	bl, _ := models.GetBlackListByCompanyID(id)
	courts, _ := models.GetCourtsByCompanyID(id)
	lawsuits, _ := models.GetLawsuitsByCompanyID(id)
	penalties, _ := models.GetPenaltiesByCompanyID(id)
	accept := this.Ctx.Request.Header.Get("accept")
	if strings.Contains(accept, "application/json") {
		this.Data["json"] = c
		this.ServeJSON()
		return
	}
	this.Data["Title"] = c.Name
	this.Data["URL"] = URL
	this.Data["Company"] = c
	this.Data["BadManagements"] = bm
	this.Data["BlackList"] = bl
	this.Data["Courts"] = courts
	this.Data["Lawsuits"] = lawsuits
	this.Data["Penalties"] = penalties
	this.Layout = "layouts/nav.html"
	this.TplName = "view.html"
	this.Render()
}

// @Title UpdateEnterprise
// @Description 修改企业信息
// @Param body body models.Company true "body for 企业 content"
// @Success 200 {int} models.Company.Id
// @Failure 403 body is empty
// @router /:id [put]
func (this *CompaniesController) Update() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less than zero or not found")
		return
	}

	c := models.Company{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &c); err != nil {
		this.CustomAbort(403, err.Error())
	}

	c.ID = id

	if err := models.UpdateCompany(&c); err != nil {
		this.CustomAbort(403, err.Error())
	}
	return
}

// @Title DeleteEnterprise
// @Description 删除企业信息
// @Param id path string true "要删除的企业id"
// @Success 204
// @Failure 403 body is empty
// @router /:id [delete]
func (this *CompaniesController) Delete() {

	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less than zero or not found")
		return
	}

	if err := models.DeleteCompany(id); err != nil {
		this.CustomAbort(403, err.Error())
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(204)

}

// @Title GETPenalty
// @Description 查询行政处罚
// @Param id path string true "企业id"
// @Param pid path string true "行政处罚id"
// @Success 200 {string}
// @Failure 403 body is empty
// @router /:id/penalties/:pid [get]
func (this *CompaniesController) GetPenalty() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if id <= 0 {
		this.Data["json"] = map[string]string{"error": "id is less zero or not found id"}
		this.ServeJSON()
		return
	}

	pid, err := this.GetInt(":pid", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if pid <= 0 {
		this.Data["json"] = map[string]string{"error": "pid is less zero or not found pid"}
		this.ServeJSON()
		return
	}

	p, err := models.GetPenaltyByIDAndCompanyID(pid, id)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	this.Data["json"] = p
	this.ServeJSON()
	return

}

// @Title AddPenalty
// @Description 添加行政处罚
// @Param id path string true "企业id"
// @Param body body models.Penalty true "body for 行政处罚"
// @Success 201 {int} models.Penalty.Id
// @Failure 403 body is empty
// @router /:id/penalty [post]
func (this *CompaniesController) AddPenalty() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less zero or not found id")
		return
	}

	p := models.Penalty{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &p)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	p.CompanyID = id

	if err := models.CreatePenalty(&p); err != nil {
		this.CustomAbort(403, err.Error())
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(201)
}

// @Title RemovePenalty
// @Description 移除行政处罚
// @Param id path string true "企业id"
// @Param pid path string true "行政处罚id"
// @Success 204
// @Failure 403 body is empty
// @router /:id/penalty/:pid [delete]
func (this *CompaniesController) RemovePenalty() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less zero or not found id")
		return
	}

	pid, err := this.GetInt(":pid", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if pid <= 0 {
		this.CustomAbort(403, "pid is less zero or not found pid")
		return
	}

	err = models.DeletePenaltyByIDAndCompanyID(pid, id)
	if err != nil {
		this.CustomAbort(403, err.Error())
	}

	this.Ctx.ResponseWriter.WriteHeader(204)

}

// @Title GETBadManagement
// @Description 查询经营异常
// @Param id path string true "企业id"
// @Param bid path string true "经营异常id"
// @Success 200 {string}
// @Failure 403 body is empty
// @router /:id/badmanagements/:bid [get]
func (this *CompaniesController) GetBadManagement() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if id <= 0 {
		this.Data["json"] = map[string]string{"error": "id is less zero or not found id"}
		this.ServeJSON()
		return
	}

	bid, err := this.GetInt(":bid", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if bid <= 0 {
		this.Data["json"] = map[string]string{"error": "bid is less zero or not found bid"}
		this.ServeJSON()
		return
	}

	bm, err := models.GetBadManagementByIDAndCompanyID(bid, id)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	this.Data["json"] = bm
	this.ServeJSON()
	return

}

// @Title AddBadManagement
// @Description 添加经营异常
// @Param id path string true
// @Param body body models.BadManagement true "body for 异常信息"
// @Success 201 {int} models.BadManagement.Id
// @Failure 403 body is empty
// @router /:id/badmanagements [post]
func (this *CompaniesController) AddBadManagement() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less zero or not found id")
		return
	}

	bm := models.BadManagement{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &bm)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	bm.CompanyID = id

	if err := models.CreateBadManagement(&bm); err != nil {
		this.CustomAbort(403, err.Error())
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(201)
}

// @Title RemoveBadManagement
// @Description 移除
// @Param id path string true
// @Param bid path string true
// @Success 204
// @Failure 403 body is empty
// @router /:id/badmanagements/:bid [delete]
func (this *CompaniesController) RemoveBadManagement() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less zero or not found id")
		return
	}

	bid, err := this.GetInt(":bid", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if bid <= 0 {
		this.CustomAbort(403, "bid is less zero or not found bid")
		return
	}

	err = models.DeleteBadManagementByIDAndCompanyID(bid, id)
	if err != nil {
		this.CustomAbort(403, err.Error())
	}

	this.Ctx.ResponseWriter.WriteHeader(204)
}

// @Title GETBlackList
// @Description 查询黑名单
// @Param id path string true "企业id"
// @Param bid path string true "黑名单id"
// @Success 200 {string}
// @Failure 403 body is empty
// @router /:id/blacklist/:bid [get]
func (this *CompaniesController) GetBlackList() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if id <= 0 {
		this.Data["json"] = map[string]string{"error": "id is less zero or not found id"}
		this.ServeJSON()
		return
	}

	bid, err := this.GetInt(":bid", -1)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	if bid <= 0 {
		this.Data["json"] = map[string]string{"error": "bid is less zero or not found bid"}
		this.ServeJSON()
		return
	}

	b, err := models.GetBlackListByIDAndCompanyID(bid, id)
	if err != nil {
		this.Data["json"] = map[string]string{"error": err.Error()}
		this.ServeJSON()
		return
	}

	this.Data["json"] = b
	this.ServeJSON()
	return

}

// @Title AddBlackList
// @Description 添加黑名单
// @Param id path string true
// @Param body body models.BlackList true "body for 黑名单"
// @Success 201 {int} models.BlackList.Id
// @Failure 403 body is empty
// @router /:id/blacklist [post]
func (this *CompaniesController) AddBlackList() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less zero or not found id")
		return
	}

	bl := models.BlackList{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &bl)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	bl.CompanyID = id

	if err := models.CreateBlackList(&bl); err != nil {
		this.CustomAbort(403, err.Error())
		return
	}
	this.Ctx.ResponseWriter.WriteHeader(201)
}

// @Title RemoveBadManagement
// @Description 移除
// @Param id path string true
// @Param bid path string true
// @Success 204
// @Failure 403 body is empty
// @router /:id/blacklist/:bid [delete]
func (this *CompaniesController) RemoveBlackList() {
	id, err := this.GetInt(":id", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if id <= 0 {
		this.CustomAbort(403, "id is less zero or not found id")
		return
	}

	bid, err := this.GetInt(":bid", -1)
	if err != nil {
		this.CustomAbort(403, err.Error())
		return
	}

	if bid <= 0 {
		this.CustomAbort(403, "bid is less zero or not found bid")
		return
	}

	err = models.DeleteBlackListByIDAndCompanyID(bid, id)
	if err != nil {
		this.CustomAbort(403, err.Error())
	}

	this.Ctx.ResponseWriter.WriteHeader(204)
}
