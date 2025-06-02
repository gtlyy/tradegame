package main

import (
	_ "tradegame/routers"
	// "myblog/models"

	// "github.com/astaxie/beego/orm"

	"github.com/astaxie/beego/logs"
	log "github.com/astaxie/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 设置日志级别和输出方式
	logs.SetLevel(logs.LevelInfo)
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/app.log","level":6,"maxlines":50000,"maxsize":1024000,"daily":true,"maxdays":10,"color":true}`)

	log.Info("Start......")
	// o := orm.NewOrm()

	// 查询
	// user1 := models.Usertb{}
	// user1.Id = 140
	// err := o.Read(&user1)
	// if err == orm.ErrNoRows {
	// 	log.Error("id=1，查询不到")
	// } else if err == orm.ErrMissPK {
	// 	log.Error("找不到主键")
	// } else {
	// 	log.Info(user1.Id, user1.Name, user1.Passwd)
	// }

	// todo: 更新，修改
	// user_update := models.Usertb{}
	// user_update.Id = 2
	// user_update.Name = "Jack"
	// user_update.Passwd = models.MD5("123456")
	// num, err := o.Update(&user_update)
	// if err != nil {
	// 	log.Println("Update fail.")
	// } else {
	// 	log.Println("Update ok, affect rows:", num)
	// }

	beego.Run()
}
