module forward-server

go 1.14

require (
	forward-core v0.0.0-00010101000000-000000000000
	github.com/astaxie/beego v1.12.2
	github.com/beego/beego/v2 v2.0.5
	github.com/dlclark/regexp2 v1.7.0
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/uber/go-torch v0.0.0-20181107071353-86f327cc820e // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	google.golang.org/appengine v1.6.7 // indirect
)

replace forward-core => ../forward-core
