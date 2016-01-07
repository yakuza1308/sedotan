package sedotan

import (
	// "bytes"
	// "fmt"
	// "github.com/eaciit/cast"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/xlsx"
	"github.com/eaciit/toolkit"
	// "reflect"
	// "regexp"
	"errors"
	// "strings"
	"time"
)

// type ViewColumn struct {
// 	Alias    string
// 	Selector string
// 	// ValueType string //-- Text, Attr, InnerHtml, OuterHtml
// 	// AttrName  string
// }

type CollectionSetting struct {
	Collection   string
	SelectColumn []*GrabColumn
	FilterCond   toolkit.M
}

type GetDatabase struct {
	dbox.ConnectionInfo
	desttype           string
	CollectionSettings map[string]*CollectionSetting

	LastExecuted time.Time

	Response []toolkit.M
}

func NewGetDatabase(host string, desttype string, connInfo *dbox.ConnectionInfo) (*GetDatabase, error) {
	g := new(GetDatabase)
	if connInfo != nil {
		g.ConnectionInfo = *connInfo
	}

	if desttype != "" {
		g.desttype = desttype
	}

	if host != "" {
		g.Host = host
	}

	if g.Host == "" || g.desttype == "" {
		return nil, errors.New("Host or Type cannot blank")
	}

	return g, nil
}

func (ds *CollectionSetting) Column(i int, column *GrabColumn) *GrabColumn {
	if i == 0 {
		ds.SelectColumn = append(ds.SelectColumn, column)
	} else if i <= len(ds.SelectColumn) {
		ds.SelectColumn[i-1] = column
	} else {
		return nil
	}
	return column
}

func (g *GetDatabase) ResultFromDatabase(dataSettingId string, out interface{}) error {

	c, e := dbox.NewConnection(g.desttype, &g.ConnectionInfo)
	if e != nil {
		return e
	}

	e = c.Connect()
	if e != nil {
		return e
	}

	defer c.Close()

	iQ := c.NewQuery()
	if g.CollectionSettings[dataSettingId].Collection != "" {
		iQ.From(g.CollectionSettings[dataSettingId].Collection)
	}

	for _, val := range g.CollectionSettings[dataSettingId].SelectColumn {
		iQ.Select(val.Selector)
	}

	//filter condition Not Yet Implemented

	csr, e := iQ.Cursor(nil)

	if e != nil {
		return e
	}
	if csr == nil {
		return e
	}
	defer csr.Close()

	ds, e := csr.Fetch(nil, 0, false)
	if e != nil {
		return e
	}

	ms := []toolkit.M{}
	for _, val := range ds.Data {
		m := toolkit.M{}
		mval := val.(toolkit.M)
		for _, column := range g.CollectionSettings[dataSettingId].SelectColumn {
			m.Set(column.Alias, "")
			if mval.Has(column.Selector) {
				m.Set(column.Alias, mval[column.Selector])
			}
		}
		ms = append(ms, m)
	}

	if edecode := toolkit.Unjson(toolkit.Jsonify(ms), out); edecode != nil {
		return edecode
	}
	return nil
}
