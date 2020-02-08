// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"testapi03/controllers"

	"github.com/astaxie/beego"
)

//func init() {
//	ns := beego.NewNamespace("/v1",
//
//		beego.NSNamespace("/conn_info",
//			beego.NSInclude(
//				&controllers.ConnInfoController{},
//			),
//		),
//
//		beego.NSNamespace("/job_info",
//			beego.NSInclude(
//				&controllers.JobInfoController{},
//			),
//		),
//
//		beego.NSNamespace("/project_image",
//			beego.NSInclude(
//				&controllers.ProjectImageController{},
//			),
//		),
//
//		beego.NSNamespace("/user_info",
//			beego.NSInclude(
//				&controllers.UserInfoController{},
//			),
//		),
//	)
//	beego.AddNamespace(ns)
//}

func init() {
	//beego.Router("/ConnInfoController", &controllers.ConnInfoController{}, )
	//beego.Router("/JobInfoController", &controllers.JobInfoController{}, )
	//beego.Router("/ProjectImageController", &controllers.ProjectImageController{},"post:Upload" )
	//beego.Router("/UserInfoController", &controllers.UserInfoController{}, )

	beego.AutoRouter(&controllers.ConnInfoController{})
	beego.AutoRouter(&controllers.JobInfoController{})
	beego.AutoRouter(&controllers.ProjectImageController{})
	beego.AutoRouter(&controllers.UserInfoController{})
}
