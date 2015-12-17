package controllers

import (
	"fmt"
	"github.com/eaciit/knot/knot.v1"
	"strings"
)

type ConfigurationController struct {
}

func (a *ConfigurationController) Default(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	return ""
}

func (a *ConfigurationController) P(k *knot.WebContext) interface{} {
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

func (a *ConfigurationController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	d := struct {
		Data string
	}{}
	e := k.GetPayload(&d)
	fmt.Println(d)
	if e != nil {
		return e.Error()
	} else {
		return d.Data
	}
}
