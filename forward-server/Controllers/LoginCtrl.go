package Controllers

import (
	"forward-core/Models"
	"forward-core/NetUtils"
	"forward-core/Utils"
	"forward-server/Controllers/BaseCtrl"
	"forward-server/Service"
	"github.com/astaxie/beego/logs"
)

type LoginCtrl struct {
	BaseCtrl.ConsoleCtrl
}

// @router /logout
func (c *LoginCtrl) Logout() {
	c.ClearUserInfo()
	c.Ctx.Redirect(302, "/login")

}

// @router /login [get]
func (c *LoginCtrl) Login() {

	c.TplName = "login.html"

}

func (c *LoginCtrl) DoLogin() {
	logs.Debug("用户登录")
	userName := c.GetString("userName")
	passWord := c.GetString("passWord")

	sysUser := Service.SysDataS.GetSysUserByName(userName)
	if sysUser == nil {
		logs.Debug("用户不存在")
		c.Ctx.Redirect(302, "/login")
		return
	}

	descryptPwd := Utils.GetMd5(passWord)
	logs.Debug("存储的密码：", sysUser.PassWord, " 输入的密码：", descryptPwd)
	if sysUser.PassWord == descryptPwd {
		if sysUser.Status == 1 {
			logs.Info("用户登录：", userName, " IP：", NetUtils.GetIP(&c.Controller))
			loginUser := new(Models.LoginUser)
			loginUser.UserId = 1
			loginUser.UserName = userName

			c.SetSession("userInfo", loginUser)
			//c.Ctx.Redirect(302, "/u/main")
			c.Data["json"] = Models.FuncResult{Code: 0, Msg: "登录成功"}
			c.ServeJSON()
		}
		if sysUser.Status == 0 {
			logs.Debug("此账号已被停用，ID：", sysUser.Id)
			//c.Ctx.Redirect(302, "/login")
			c.Data["json"] = Models.FuncResult{Code: 1, Msg: "您的账号已被停用"}
			c.ServeJSON()
		}
	} else {
		logs.Debug("用户登录失败")
		// c.Ctx.Redirect(302, "/login")
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "登录失败"}
		c.ServeJSON()
	}

}
