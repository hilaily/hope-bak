package controllers

import (
	"encoding/json"
	"hope/app/models"
	"hope/app/routes"
	"hope/app/support"
	"strconv"
	"time"

	"github.com/revel/revel"
)

//Login controller
type Login struct {
	*revel.Controller
}

//SignIn page.
func (l Login) SignIn() revel.Result {
	return l.Render()
}

//handle Sign
func (l Login) SignInHandler(name, passwd string) revel.Result {

	l.Validation.Required(name).Message("username can't be null.")
	l.Validation.Required(passwd).Message("passwd can't be null.")

	if l.Validation.HasErrors() {
		l.Validation.Keep()
		l.FlashParams()
		return l.Redirect(routes.Login.SignIn())
	}

	model := &models.Admin{Name: name, Passwd: passwd}
	revel.WARN.Printf("用户[name:%s; pass:%s]尝试登陆", name, passwd)
	admin, err := model.SignIn(l.Request)
	if err != "" {
		revel.ERROR.Println(err)
		return l.Redirect(routes.Login.SignIn())
	}
	revel.WARN.Printf("用户[name:%s; id:%d]登陆成功", admin.Name, admin.Id)

	//put admin id in seesion
	l.Session["UID"] = strconv.Itoa(int(admin.Id))
	//set admin info in cache, time out time.Minute * 30
	support.MCache.Delete(support.SPY_ADMIN_INFO + strconv.Itoa(int(admin.Id)))
	data, _ := json.Marshal(&admin)
	support.MCache.Set(support.SPY_ADMIN_INFO+strconv.Itoa(int(admin.Id)), string(data), time.Minute*30)

	return l.Redirect(routes.Admin.Main())
}

//SignUp page.
func (l Login) SignUp() revel.Result {
	return l.Render()
}

// SignOut
func (l Login) SignOut() revel.Result {
	uid := l.Session["UID"]
	if uid != "" {
		delete(l.Session, "UID")
	}
	return l.Redirect(routes.Login.SignIn())
}

//handle sign up.
func (l Login) SignUpHandler(name, email, passwd string) revel.Result {

	l.Validation.Required(name).Message("username can't be null.")
	l.Validation.Required(email).Message("email can't be null.")
	l.Validation.Required(passwd).Message("passwd can't be null.")

	if l.Validation.HasErrors() {
		l.Validation.Keep()
		l.FlashParams()
		return l.Redirect(routes.Login.SignUp())
	}

	model := &models.Admin{Name: name, Email: email, Passwd: passwd}
	has, err := model.New()

	if err != "" && has <= 0 {
		revel.ERROR.Println(err)
		l.Flash.Error("msg", err)
		return l.Redirect(routes.Login.SignUp())
	}

	l.Flash.Success("msg", "success")
	return l.Redirect(routes.Login.SignIn())
}
