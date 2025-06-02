package routers

import (
	"tradegame/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// beego.Router("/", &controllers.MainController{})
	beego.Router("/", &controllers.TradeController{})

	beego.Router("/register", &controllers.RegController{})
	beego.Router("/check-username", &controllers.UserController{}, "post:CheckUsername")

	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/info", &controllers.InfoController{})

	beego.Router("/qid", &controllers.QidController{})

	beego.Router("/trade", &controllers.TradeController{})

}
