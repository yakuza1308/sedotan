package sedotan

import (
	"fmt"
	//gq "github.com/PuerkitoBio/goquery"
	"github.com/eaciit/toolkit"
	"net/http"
	"time"
)

type GrabConfig struct {
	Data         toolkit.M
	AuthType     string
	AuthUserId   string
	AuthPassword string
}

type Grabber struct {
	URL      string
	CallType string
	Config   *GrabConfig

	LastExecuted time.Time

	bodyByte []byte
	Response *http.Response
}

func NewGrabber(url string, calltype string, config *GrabConfig) *Grabber {
	g := new(Grabber)
	g.URL = url
	g.CallType = calltype
	if config == nil {
		config = new(GrabConfig)
	}
	g.Config = config
	g.bodyByte = []byte{}
	return g
}

func (g *Grabber) Data() interface{} {
	return nil
}

func (g *Grabber) DataByte() []byte {
	d := g.Data()
	if toolkit.IsValid(d) {
		return toolkit.Jsonify(d)
	}
	return []byte{}
}

func (g *Grabber) Grab(parm toolkit.M) error {
	r, e := toolkit.HttpCall(g.URL, g.CallType, g.DataByte(), nil)
	errorTxt := ""
	if e != nil {
		errorTxt = e.Error()
	} else if r.StatusCode != 200 {
		errorTxt = r.Status
	}
	if errorTxt != "" {
		return fmt.Errorf("Unable to grab %s. %s", g.URL, errorTxt)
	}

	g.Response = r
	g.bodyByte = toolkit.HttpContent(r)
	return nil
}

func (g *Grabber) ResultString() string {
	if g.Response == nil {
		return ""
	}

	return string(g.bodyByte)
}

/*
func (g *Grabber) ResultPart() (string, error){
	s := g.ResultString()
	doc, e := gq.NewDocument(s)
	if e!=nil {
		return "", e.Error()
	}

	f := doc.Find(".content")
	retun f.First().Text(), nil
}
*/
