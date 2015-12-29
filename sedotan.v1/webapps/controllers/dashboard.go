package controllers

import (
	// "fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/json"
	"github.com/eaciit/knot/knot.v1"
	// sdt "github.com/eaciit/sedotan/sedotan.v1"
	"github.com/eaciit/sedotan/sedotan.v1/webapps/modules"
	// "github.com/eaciit/toolkit"
	"strings"
)

var (
	filename = wd + "data\\config.json"
)

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

func (a *DashboardController) Griddashboard(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

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
	csr, e := c.NewQuery().Select("nameid", "url", "grabinterval").Cursor(nil)
	defer csr.Close()
	ds, e := csr.Fetch(nil, 0, false)

	return ds.Data
}

func (a *DashboardController) Startservice(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	t := struct {
		NameId string
	}{}
	e := k.GetPayload(&t)
	if e != nil {
		return e.Error()
	}

	ds, _ := Getquery(t.NameId)
	_, isRun := modules.Process(ds)

	return isRun
}

func (a *DashboardController) Stopservice(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	t := struct {
		NameId string
	}{}
	e := k.GetPayload(&t)
	if e != nil {
		return e.Error()
	}

	ds, _ := Getquery(t.NameId)
	_, isRun := modules.StopProcess(ds)
	// var grabName = map[string]interface{}{}
	// grabName["name"] = name
	// grabName["stat"] = isRun

	return isRun
}

func (a *DashboardController) Stat(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	t := struct {
		NameId   string
		BtnClick string
	}{}
	e := k.GetPayload(&t)
	if e != nil {
		return e.Error()
	}

	ds, _ := Getquery(t.NameId)
	grabStatus := modules.CheckStat(ds)

	return grabStatus
}

func Getquery(nameid string) ([]interface{}, string) {
	ci := &dbox.ConnectionInfo{filename, "", "", "", nil}
	c, e := dbox.NewConnection("json", ci)
	if e != nil {
		return nil, e.Error()
	}

	e = c.Connect()
	if e != nil {
		return nil, e.Error()
	}
	defer c.Close()

	csr, e := c.NewQuery().Where(dbox.Eq("nameid", nameid)).Cursor(nil)

	ds, e := csr.Fetch(nil, 0, false)
	return ds.Data, ""
}
