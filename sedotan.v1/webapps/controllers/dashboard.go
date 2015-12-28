package controllers

import (
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/json"
	"github.com/eaciit/knot/knot.v1"
	// "github.com/eaciit/sedotan/sedotan.v1/webapps/modules"
	// "github.com/eaciit/toolkit"
	"strings"
)

// var (
// 	wd = func() string {
// 		d, _ := os.Getwd()
// 		return d + "/../"
// 	}()
// )

// func init() {
// 	app := knot.NewApp("sedotan")
// 	app.ViewsPath = wd + "views/"
// 	app.Register(&AppController{})
// 	app.Static("static", wd+"assets")
// 	app.LayoutTemplate = "_layout.html"
// 	knot.RegisterApp(app)
// }

type DashboardController struct {
}

func (a *DashboardController) Default(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	return ""
}

func (a *DashboardController) P(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	vn := ""
	qs := k.Request.RequestURI
	if qs != "" {
		qss := strings.Split(qs, "?")
		if len(qss) > 1 {
			vn = strings.Split(qss[1], "&")[0]
			if strings.HasSuffix(vn, ".html") == false {
				vn += ".html"
			}
		}
		k.Config.ViewName = vn
	}
	return ""
}

func (a *DashboardController) Getconfig(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	filename := wd + "data\\config.json"
	ci := &dbox.ConnectionInfo{filename, "", "", "", nil}

	c, e := dbox.NewConnection("json", ci)
	if e != nil {
		return e
	}

	e = c.Connect()
	if e != nil {
		return e
	}
	defer c.Close()

	csr, e := c.NewQuery().Select("*").Cursor(nil)
	defer csr.Close()

	ds, e := csr.Fetch(nil, 0, false)

	// s := modules.Grabget(ds.Data)
	for _, v := range ds.Data {
		// fmt.Printf("datas:%v\n", v.(map[string]interface{})["nameid"])
		fmt.Println(v.(map[string]interface{})["grabinterval"])
	}

	// t := struct {
	// 	Name string
	// }{}
	// e := k.GetPayload(&t)
	// if e != nil {
	// 	return e.Error()
	// } else {
	// 	return "Hi " + t.Name
	// }
	return ds.Data
}

func (a *DashboardController) Stat(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// e = xGrabService.StartService()
	// if e != nil {
	// 	t.Errorf("Error Found : ", e)
	// } else {

	// 	for i := 0; i < 100; i++ {
	// 		fmt.Printf(".")
	// 		time.Sleep(3000 * time.Millisecond)
	// 	}

	// 	e = xGrabService.StopService()
	// 	if e != nil {
	// 		t.Errorf("Error Found : ", e)
	// 	}
	// }
	return "run"
}
