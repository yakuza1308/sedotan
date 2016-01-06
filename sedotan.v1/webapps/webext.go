package webext

import (
	"github.com/eaciit/knot/knot.v1"
	. "github.com/eaciit/sedotan/sedotan.v1/webapps/controllers"
	. "github.com/eaciit/sedotan/sedotan.v1/webapps/modules"
	"os"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/../"
	}()
)

func init() {
	app := knot.NewApp("sedotan")
	app.ViewsPath = wd + "views/"
	app.Controllers()
	// app.Register(&AppController{})
	app.Register(&DashboardController{})
	app.Register(new(ResultController))
	app.Register(new(ConfigurationController))
	app.Register(new(GrabModule))
	app.Static("static", wd+"assets")
	app.LayoutTemplate = "_layout.html"
	knot.RegisterApp(app)

}
