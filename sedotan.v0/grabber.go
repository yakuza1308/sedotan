package sedotan

import (
	"bytes"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"github.com/eaciit/toolkit"
	"net/http"
	"strings"
	"time"
)

type GrabColumn struct {
	Alias     string
	Selector  string
	ValueType string //-- Text, Attr, InnerHtml, OuterHtml
	AttrName  string
}

type Config struct {
	Data         toolkit.M
	URL          string
	CallType     string
	DataSettings map[string]*DataSetting
	FormValues   toolkit.M

	AuthType     string
	AuthUserId   string
	AuthPassword string
}

type DataSetting struct {
	RowSelector    string
	RowDeleteCond  toolkit.M
	ColumnSettings []*GrabColumn
}

type Grabber struct {
	Config

	LastExecuted time.Time

	bodyByte []byte
	Response *http.Response
}

func NewGrabber(url string, calltype string, config *Config) *Grabber {
	g := new(Grabber)
	if config != nil {
		g.Config = *config
	}

	if url != "" {
		g.URL = url
	}

	if calltype != "" {
		g.CallType = calltype
	}

	return g
}

func (c *Config) SetFormValues(parm toolkit.M) {
	c.FormValues = toolkit.M{}.Set("formvalues", parm)
}

func (ds *DataSetting) Column(i int, column *GrabColumn) *GrabColumn {
	if i == 0 {
		ds.ColumnSettings = append(ds.ColumnSettings, column)
	} else if i <= len(ds.ColumnSettings) {
		ds.ColumnSettings[i-1] = column
	} else {
		return nil
	}
	return column
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

	r, e := toolkit.HttpCall(g.URL, g.CallType, g.DataByte(), g.Config.FormValues)
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

func (g *Grabber) ResultFromHtml(dataSettingId string, out interface{}) error {

	reader := bytes.NewReader(g.bodyByte)
	doc, e := gq.NewDocumentFromReader(reader)
	if e != nil {
		return e
	}

	ms := []toolkit.M{}
	records := doc.Find(g.Config.DataSettings[dataSettingId].RowSelector)
	recordCount := records.Length()

	for i := 0; i < recordCount; i++ {
		record := records.Eq(i)
		m := toolkit.M{}
		for cindex, c := range g.Config.DataSettings[dataSettingId].ColumnSettings {
			columnId := fmt.Sprintf("%s", cindex)
			if c.Alias != "" {
				columnId = c.Alias
			}
			sel := record.Find(c.Selector)
			var value interface{}
			valuetype := strings.ToLower(c.ValueType)
			if valuetype == "attr" {
				value, _ = sel.Attr(c.AttrName)
			} else if valuetype == "html" {
				value, _ = sel.Html()
			} else {
				value = sel.Text()
			}
			value = strings.TrimSpace(fmt.Sprintf("%s", value))
			m.Set(columnId, value)
		}

		if !(g.Config.DataSettings[dataSettingId].getCondition(m)) {
			ms = append(ms, m)
		}
	}
	if edecode := toolkit.Unjson(toolkit.Jsonify(ms), out); edecode != nil {
		return edecode
	}
	return nil
}

func (ds *DataSetting) getCondition(dataCheck toolkit.M) bool {
	resBool := true

	if len(ds.RowDeleteCond) > 0 {
		resBool = foundCondition(dataCheck, ds.RowDeleteCond)
	}

	return resBool
}

func foundCondition(dataCheck toolkit.M, cond toolkit.M) bool {
	resBool := true

	for key, val := range cond {
		if key == "$and" || key == "$or" {
			for i, sVal := range val.([]interface{}) {
				rVal := sVal.(map[string]interface{})
				mVal := toolkit.M{}
				for rKey, mapVal := range rVal {
					mVal.Set(rKey, mapVal)
				}

				xResBool := foundCondition(dataCheck, mVal)
				if key == "$and" {
					resBool = resBool && xResBool
				} else {
					if i == 0 {
						resBool = xResBool
					} else {
						resBool = resBool || xResBool
					}
				}
			}
		} else if val != dataCheck.Get(key, "").(string) {
			resBool = false
		}
	}

	return resBool
}
