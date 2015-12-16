package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	"strings"
)

type AppController struct {
}

func (a *AppController) Default(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	return ""
}

func (a *AppController) P(k *knot.WebContext) interface{} {
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

func (a *AppController) Process(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	t := struct {
		Nama string
	}{}
	e := k.GetPayload(&t)
	if e != nil {
		return e.Error()
	} else {
		return "Hi " + t.Nama
	}
}

func (a *AppController) ProcessForm(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	var model struct {
		Nama string
		Umur int
	}
	e := k.GetForms(&model)
	if e != nil {
		model.Nama = e.Error()
	}
	return model
}
