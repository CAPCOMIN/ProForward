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
	beego.Router("/performance/cpu", &Controllers.PerformanceController{}, "get:PerfCPU")
	beego.Router("/performance/heap", &Controllers.PerformanceController{}, "get:PerfHeap")
	beego.Router("/performance/tbgc", &Controllers.PerformanceController{}, "get:PerfThreadBlockGC")
	beego.Router("/performance/monitorcpu", &Controllers.PerformanceController{}, "get:MonitorCPU")
	beego.Router("/performance/domonitorcpu", &Controllers.PerformanceController{}, "post:DoMonitorCPU")
	beego.Router("/performance/mem", &Controllers.PerformanceController{}, "get:PerfMemInfo")
	beego.Router("/performance/getmem", &Controllers.PerformanceController{}, "get:DoPerfMemInfo")
	beego.Router("/performance/pprof", &Controllers.PerformanceController{}, "get:PprofWebAnalysis")

	beego.Router("/u/EditForward", &Controllers.ForwardCtrl{}, "get:EditForward")
	beego.Router("/u/DoEditForward", &Controllers.ForwardCtrl{}, "post:DoEditForward")

	beego.Router("/u/Dosaveforward", &Controllers.ForwardCtrl{}, "post:DoSaveForward")

	beego.Router("/u/Dologins", &Controllers.LoginCtrl{}, "post:DoLogin")

	beego.Router("/u/about", &Controllers.ForwardCtrl{}, "get:AboutForward")
	//beego.SetStaticPath("/swagger", "swagger")

}
