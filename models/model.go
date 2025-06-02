package models

import (
	"crypto/md5"
	"fmt"

	log "github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 数据模型。备注：需要与数据库中的表名和字段名保持一致。（貌似大小写可以忽略）
type Usertb struct {
	Id     int
	Name   string
	Passwd string
}

// 导入包，自动运行 init() 函数。
func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "tgame:tc_3469Uk@tcp(localhost:3306)/webdb?charset=utf8")
	orm.RegisterModel(new(Usertb))
}

// 传入的数据不一样，那么MD5后的32位长度的数据肯定会不一样
func MD5(str string) string {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5str
}

// 检查用户名是否已存在
func CheckUsernameExists(username string) bool {
	log.Info("In model.go: CheckUsernameExists().")
	o := orm.NewOrm()
	exist := o.QueryTable("usertb").Filter("name", username).Exist()
	log.Info("exist=", exist)
	return exist
}
