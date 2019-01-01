package controllers

import (
	"fmt"
	"hope/app/models"
	"hope/app/support"

	"github.com/revel/revel"
)

type DBParams struct {
	Db_host   string
	Db_user   string
	Db_pass   string
	Db_port   string
	Db_name   string
	Db_prefix string
	Driver    string
}

type AdminParams struct {
	Admin_user  string
	Admin_pass  string
	Admin_email string
}

type Install struct {
	*revel.Controller
}

func (i *Install) Index() revel.Result {
	return i.Render()
}

func (i *Install) HandleInstall() revel.Result {
	return nil
}

// AddAdmin to add a admin user when installing
// 用于在安装的时候添加一个管理员账号
func (i *Install) AddAdmin() revel.Result {
	params := new(AdminParams)
	i.Params.Bind(params, "info")
	admin := &models.Admin{Name: params.Admin_user, Nickname: params.Admin_user, Passwd: params.Admin_pass, Email: params.Admin_email, RoleId: models.ADMIN_SUPER}
	id, msg := admin.New()

	if id <= 0 {
		return i.RenderJSON(&ResultJson{Success: false, Msg: msg, Data: ""})
	}
	return i.RenderJSON(&ResultJson{Success: true})
}

func (i *Install) AddDB() revel.Result {
	params := new(DBParams)
	i.Params.Bind(params, "info")
	params.Driver = "mysql"
	err := i.checkDB(params)
	if err != nil {
		msg := "连接数据库失败：" + err.Error()
		revel.ERROR.Println(msg)
		return i.RenderJSON(&ResultJson{Success: false, Msg: msg, Data: ""})
	}
	revel.INFO.Println("开始同步数据库...")
	err = models.SyncDB()
	if err != nil {
		msg := "同步数据库失败：" + err.Error()
		revel.ERROR.Println(msg)
		return i.RenderJSON(&ResultJson{Success: false, Msg: msg, Data: ""})
	}
	revel.INFO.Println("同步数据库完成...")
	admin := new(models.Admin)
	adminCount := admin.AdminSuperUserCount()
	i.finishInstall()
	return i.RenderJSON(&ResultJson{Success: true, Data: adminCount})
}

func (i *Install) checkDB(info *DBParams) error {
	fmt.Println("db info: ", info)
	err := support.TestXorm(info.Driver, info.Db_user, info.Db_pass, info.Db_host, info.Db_port, info.Db_name, info.Db_prefix)
	if err != nil {
		return err
	}
	err = support.AddDB(info.Db_host, info.Db_port, info.Db_user, info.Db_pass, info.Db_name, info.Db_prefix, info.Driver)
	if err != nil {
		return err
	}
	return nil
}

func (i *Install) finishInstall() {
	support.IsInstalled = true
	support.FinishInstall()
}
