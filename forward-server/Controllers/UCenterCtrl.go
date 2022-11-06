package Controllers

import (
	"forward-core/Models"
	"forward-core/Utils"
	"forward-server/Controllers/BaseCtrl"
	"forward-server/Service"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strconv"
	"time"
)

type UCenterCtrl struct {
	BaseCtrl.ConsoleCtrl
}

// @router /u/main [get]
func (c *UCenterCtrl) Main() {

	c.Data["loginUserName"] = c.LoginUser.UserName
	c.Layout = "ucenter/layout.html"
	c.TplName = "ucenter/main.html"

}

// @router /u/index [get]
func (c *UCenterCtrl) Index() {

	c.Data["runtime_NumCPU"] = runtime.NumCPU()
	c.Data["runtime_GOOS"] = runtime.GOOS
	c.Data["runtime_GOARCH"] = runtime.GOARCH
	c.Data["runtime_NumGoroutine"] = runtime.NumGoroutine()
	c.Data["server_Time"] = time.Now()
	//c.Data["runtime_NumGoroutine"] = runtime.

	c.TplName = "ucenter/index.html"
}

// @router /u/getServerTime [post]
func (c *UCenterCtrl) GetServerTime() {

	c.Data["json"] = Utils.GetCurrentTime()

	c.ServeJSON()

}

// @router /u/userManage [get]
func (c *UCenterCtrl) UserManage() {

	userList := Service.SysDataS.GetAllSysUser()

	c.Data["userList"] = userList

	//logs.Warn("UserManage", userList[0].UserName)

	c.TplName = "ucenter/userManage.html"
}

func (c *UCenterCtrl) AddUser() {
	var newUser Models.SysUser

	newUser.UserName = c.GetString("UserName")
	newUser.PassWord = Utils.GetMd5(c.GetString("PassWord"))
	newUser.CreateTime = time.Now()
	Switch := c.GetString("Status")

	if Switch == "on" {
		newUser.Status = 0
	} else {
		newUser.Status = 1
	}
	logs.Warn("AddUser", newUser.UserName, newUser.PassWord, newUser.Status, newUser.CreateTime)

	id, err := Service.SysDataS.AddOneUser(&newUser)
	if err != nil {
		logs.Error("AddUser", err)
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "添加账户失败，" + err.Error()}
	} else {
		c.Data["json"] = Models.FuncResult{Code: 0, Msg: "已成功添加账户，ID:" + strconv.FormatInt(id, 10)}
	}
	c.ServeJSON()
}

func (c *UCenterCtrl) BanUser() {
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("BanUser GetInt")
	}
	err2 := Service.SysDataS.ChangeSysUserStatus(id, 0)
	if err2 != nil || err != nil {
		logs.Error("BanUser2", err2)
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "停用账户失败，" + err.Error()}
	} else {
		c.Data["json"] = Models.FuncResult{Code: 0, Msg: "已成功停用账户，ID:"}
	}
	c.ServeJSON()
}
func (c *UCenterCtrl) EnableUser() {
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("EnableUser GetInt")
	}
	err2 := Service.SysDataS.ChangeSysUserStatus(id, 1)
	if err2 != nil || err != nil {
		logs.Error("EnableUser2", err2)
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "启用账户失败，" + err.Error()}
	} else {
		c.Data["json"] = Models.FuncResult{Code: 0, Msg: "已成功启用账户，ID:" + strconv.Itoa(id)}
	}
	c.ServeJSON()
}

func (c *UCenterCtrl) DeleteOneUser() {
	id, err := c.GetInt("id")
	if err != nil {
		logs.Error("DeleteOneUser GetInt")
	}
	err2 := Service.SysDataS.DelOneSysUsers(id)
	if err2 != nil || err != nil {
		logs.Error("DeleteOneUser", err2)
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "删除账户失败，" + err.Error()}
	} else {
		c.Data["json"] = Models.FuncResult{Code: 0, Msg: "已成功删除账户，ID:" + strconv.Itoa(id)}
		// c.Ctx.Redirect(302, "/u/main")
	}
	c.ServeJSON()
}

// @router /u/changePwd [get]
func (c *UCenterCtrl) ChangePwd() {

	c.TplName = "ucenter/changePwd.html"
}

// @router /u/doChangePwd [post]
func (c *UCenterCtrl) DoChangePwd() {
	userInfo := c.GetUserInfo()

	passWord := c.GetString("passWord")
	passWord2 := c.GetString("passWord2")

	if Utils.IsEmpty(passWord) {
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "密码不能为空"}
		c.ServeJSON()
		return
	}

	if passWord != passWord2 {
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: "两次输入的密码不一致"}
		c.ServeJSON()
		return
	}

	err := Service.SysDataS.ChangeUserPwd(userInfo.UserId, passWord)
	if err == nil {
		c.Data["json"] = Models.FuncResult{Code: 0, Msg: "密码修改成功"}
	} else {
		c.Data["json"] = Models.FuncResult{Code: 1, Msg: err.Error()}
	}

	c.ServeJSON()

}
