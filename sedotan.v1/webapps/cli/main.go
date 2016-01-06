package main

import (
	"github.com/eaciit/knot/knot.v1"
	_ "github.com/eaciit/sedotan/sedotan.v1/webapps"
)

func main() {
	app := knot.GetApp("sedotan")
	if app == nil {
		return
	}
	knot.StartApp(app, "localhost:1308")

}
