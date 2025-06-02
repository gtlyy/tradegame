package controllers

import (
	log "github.com/astaxie/beego/logs"
	// "myblog/models"

	// "github.com/astaxie/beego/orm"

	beego "github.com/beego/beego/v2/server/web"
)

// 用户界面控制器
type InfoController struct {
	beego.Controller
}

//
func (c *InfoController) Get() {
	log.Info("In InfoController: Get().")
	userName := c.Ctx.GetCookie("userName")
	c.Data["userName"] = userName
	c.TplName = "info.html"
}
