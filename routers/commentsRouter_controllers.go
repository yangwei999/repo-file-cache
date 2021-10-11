package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/opensourceways/repo-file-cache/controllers:FileController"] = append(beego.GlobalControllerRouter["github.com/opensourceways/repo-file-cache/controllers:FileController"],
		beego.ControllerComments{
			Method:           "Set",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/opensourceways/repo-file-cache/controllers:FileController"] = append(beego.GlobalControllerRouter["github.com/opensourceways/repo-file-cache/controllers:FileController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           "/:platform/:org/:repo/:branch/:filename",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
