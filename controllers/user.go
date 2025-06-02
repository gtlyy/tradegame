package controllers

import (
	// "myblog/models"

	"tradegame/models"

	log "github.com/astaxie/beego/logs"

	// "github.com/astaxie/beego/orm"

	beego "github.com/beego/beego/v2/server/web"
)

// 用户控制器
type UserController struct {
	beego.Controller
}

// 检查用户名
func (c *UserController) CheckUsername() {
	log.Info("In user.go: CheckUsername().")
	username := c.GetString("username")

	// 检查用户名是否已存在
	log.Info("Checking if a Username Already Exists: ", username)
	if models.CheckUsernameExists(username) {
		c.Data["json"] = map[string]interface{}{
			"available": false, // 用户已存在，不可用。
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"available": true,
		}
	}

	c.ServeJSON()
}
