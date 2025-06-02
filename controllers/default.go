package controllers

import (
	"tradegame/models"

	log "github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

// Reg ============================================================
type RegController struct {
	beego.Controller
}

func (c *RegController) Get() {
	c.TplName = "register.html"
}

func (c *RegController) Post() {
	username := c.GetString("username")
	password := c.GetString("password")
	log.Info("In default.go: Post(): register:", username, password)
	// 存入数据库：
	log.Info("In default.go: Post(): Store into db.")
	o := orm.NewOrm()
	user := new(models.Usertb)
	user.Name = username
	user.Passwd = password
	n, err := o.Insert(user)
	if err != nil {
		log.Info("Error: Insert()", err)
	}
	log.Info("In defalut.go: Reg Post(): ", n)
	o.Commit() // 提交事务
	c.Data["json"] = map[string]interface{}{
		"id": n,
	}
	c.ServeJSON()
}

// Login ==========================================================
type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	log.Info("In default.go: Login: Get().")
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	log.Info("In default.go: Login: Post().")
	username := c.GetString("username")
	password := c.GetString("password")
	// log.Info(username, password)

	o := orm.NewOrm()
	user := new(models.Usertb)
	user.Name = username
	user.Passwd = password
	err := o.QueryTable("usertb").Filter("name", username).Filter("passwd", password).One(user)
	if err == nil {
		ip := c.Ctx.Input.IP()
		log.Info("Login ok: Name, IP=", user.Name, ip)
		c.Data["json"] = map[string]interface{}{
			"result":   true,
			"redirect": "/info",
		}
		// 登录验证成功，将用户ID存储到 Cookie 中
		c.Ctx.SetCookie("userID", "123", 3600)      // 参数分别为 Cookie 名称、值和过期时间
		c.Ctx.SetCookie("userName", username, 3600) // 参数分别为 Cookie 名称、值和过期时间
		c.ServeJSON()
	} else {
		c.Data["json"] = map[string]interface{}{
			"result": false,
			"err":    "Invalid username or password",
		}
		c.ServeJSON()
	}

}
