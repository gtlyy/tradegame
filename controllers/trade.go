package controllers

import (
	log "github.com/astaxie/beego/logs"

	beego "github.com/beego/beego/v2/server/web"
)

// 用户界面控制器
type TradeController struct {
	beego.Controller
}

//
func (c *TradeController) Get() {
	log.Info("In TradeController: Get().")
	ip := c.Ctx.Input.IP()
	log.Info("Welcome IP:", ip)
	u := c.Ctx.GetCookie("userID")
	log.Info("UserID:", u)
	if u != "123" {
		// c.Redirect("/login", 302)
		// return
		c.TplName = "nologin.html"
		return
	}

	c.TplName = "tradegame.html"
}
