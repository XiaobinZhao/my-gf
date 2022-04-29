package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	swaggerUIDefaultURL = `https://petstore.swagger.io/v2/swagger.json`
)

func SwaggerUI(r *ghttp.Request) {
	var (
		indexFileName = `index.html`
	)
	if r.StaticFile != nil && r.StaticFile.File != nil && gfile.Basename(r.StaticFile.File.Name()) == indexFileName {
		//if gfile.Basename(r.URL.Path) != indexFileName && r.originUrlPath[len(r.originUrlPath)-1] != '/' {
		//	r.Response.Header().Set("Location", r.originUrlPath+"/")
		//	r.Response.WriteHeader(http.StatusMovedPermanently)
		//	r.ExitAll()
		//}
		r.Response.Write(gstr.Replace(
			string(r.StaticFile.File.Content()),
			swaggerUIDefaultURL,
			"/swagger",
		))
		r.ExitAll()
	}
}
