package routers

import (
	"forward-server/Controllers"
	"github.com/astaxie/beego"
)

func init() {

	//
	beego.Include(&Controllers.DefaultCtrl{})
	beego.Include(&Controllers.LoginCtrl{})
	beego.Include(&Controllers.UCenterCtrl{})
	beego.Include(&Controllers.ForwardCtrl{})

	api_ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&Controllers.RestApiCtrl{},
			),
		),
	)

	beego.AddNamespace(api_ns)
	beego.Router("/manageUser/ban", &Controllers.UCenterCtrl{}, "get:BanUser")
	beego.Router("/manageUser/enable", &Controllers.UCenterCtrl{}, "get:EnableUser")

	beego.Router("/manageUser/del", &Controllers.UCenterCtrl{}, "get:DeleteOneUser")

	beego.Router("/manageUser/add", &Controllers.UCenterCtrl{}, "post:AddUser")

	beego.Router("/manageUser/edit", &Controllers.UCenterCtrl{}, "get:EditUser")

	beego.Router("/manageUser/doEdit", &Controllers.UCenterCtrl{}, "post:DoEditUser")

	beego.Router("/performance/goroutine", &Controllers.PerformanceController{}, "get:PerfGoroutine")
	
	beego.Router("/u/EditForward", &Controllers.ForwardCtrl{}, "get:EditForward")
	
	beego.Router("/u/DoEditForward", &Controllers.ForwardCtrl{}, "post:DoEditForward")
	
	//beego.SetStaticPath("/swagger", "swagger")

}
