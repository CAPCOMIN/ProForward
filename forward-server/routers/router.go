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

}
