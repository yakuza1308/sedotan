package sedotan

import (
	"bytes"
	"fmt"
	"github.com/eaciit/cast"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/dbox/dbc/xlsx"
	"github.com/eaciit/toolkit"
	"reflect"
	// "regexp"
	"errors"
	"strings"
	"time"
)

type ViewColumn struct {
	Alias    string
	Selector string
	// ValueType string //-- Text, Attr, InnerHtml, OuterHtml
	// AttrName  string
}

type CollectionSetting struct {
	Collection   string
	SelectColumn []*ViewColumn
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
	if config != nil {
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

func (ds *CollectionSetting) Column(i int, column *ViewColumn) *ViewColumn {
	if i == 0 {
		ds.SelectColumn = append(ds.SelectColumn, column)
	} else if i <= len(ds.SelectColumn) {
		ds.SelectColumn[i-1] = column
	} else {
		return nil
	}
	return column
}

func (g *GetDatabase) ResultFromDB(dataSettingId string, out interface{}) error {

	c, e := dbox.NewConnection(g.desttype, g.ConnectionInfo)
	if e != nil {
		return e
	}

	e = c.Connect()
	if e != nil {
		return e
	}

	c.Close()

	ms := []toolkit.M{}

	// for i := 0; i < recordCount; i++ {
	// 	record := records.Eq(i)
	// 	m := toolkit.M{}
	// 	for cindex, c := range g.Config.DataSettings[dataSettingId].ColumnSettings {
	// 		columnId := fmt.Sprintf("%s", cindex)
	// 		if c.Alias != "" {
	// 			columnId = c.Alias
	// 		}
	// 		sel := record.Find(c.Selector)
	// 		var value interface{}
	// 		valuetype := strings.ToLower(c.ValueType)
	// 		if valuetype == "attr" {
	// 			value, _ = sel.Attr(c.AttrName)
	// 		} else if valuetype == "html" {
	// 			value, _ = sel.Html()
	// 		} else {
	// 			value = sel.Text()
	// 		}
	// 		value = strings.TrimSpace(fmt.Sprintf("%s", value))
	// 		m.Set(columnId, value)
	// 	}

	// 	// if !(g.Config.DataSettings[dataSettingId].getDeleteCondition(m)) {
	// 	ms = append(ms, m)
	// 	// }
	// }
	if edecode := toolkit.Unjson(toolkit.Jsonify(ms), out); edecode != nil {
		return edecode
	}
	return nil
}

// func (ds *DataSetting) getDeleteCondition(dataCheck toolkit.M) bool {
// 	resBool := false

// 	if len(ds.RowDeleteCond) > 0 {
// 		resBool = foundCondition(dataCheck, ds.RowDeleteCond)
// 	}

// 	if !resBool && len(ds.RowIncludeCond) > 0 {
// 		resBool = true
// 		iResBool := true
// 		iResBool = foundCondition(dataCheck, ds.RowIncludeCond)
// 		if iResBool {
// 			resBool = false
// 		}
// 	}

// 	return resBool
// }

// func foundCondition(dataCheck toolkit.M, cond toolkit.M) bool {
// 	resBool := true

// 	for key, val := range cond {
// 		if key == "$and" || key == "$or" {
// 			for i, sVal := range val.([]interface{}) {
// 				rVal := sVal.(map[string]interface{})
// 				mVal := toolkit.M{}
// 				for rKey, mapVal := range rVal {
// 					mVal.Set(rKey, mapVal)
// 				}

// 				xResBool := foundCondition(dataCheck, mVal)
// 				if key == "$and" {
// 					resBool = resBool && xResBool
// 				} else {
// 					if i == 0 {
// 						resBool = xResBool
// 					} else {
// 						resBool = resBool || xResBool
// 					}
// 				}
// 			}
// 		} else {
// 			if reflect.ValueOf(val).Kind() == reflect.Map {
// 				mVal := val.(map[string]interface{})
// 				tomVal, _ := toolkit.ToM(mVal)
// 				switch {
// 				case tomVal.Has("$ne"):
// 					if tomVal["$ne"].(string) == dataCheck.Get(key, "").(string) {
// 						resBool = false
// 					}
// 				case tomVal.Has("$regex"):
// 					resBool, _ = regexp.MatchString(tomVal["$regex"].(string), dataCheck.Get(key, "").(string))
// 				case tomVal.Has("$gt"):
// 					if tomVal["$gt"].(string) >= dataCheck.Get(key, "").(string) {
// 						resBool = false
// 					}
// 				case tomVal.Has("$gte"):
// 					if tomVal["$gte"].(string) > dataCheck.Get(key, "").(string) {
// 						resBool = false
// 					}
// 				case tomVal.Has("$lt"):
// 					if tomVal["$lt"].(string) <= dataCheck.Get(key, "").(string) {
// 						resBool = false
// 					}
// 				case tomVal.Has("$lte"):
// 					if tomVal["$lte"].(string) < dataCheck.Get(key, "").(string) {
// 						resBool = false
// 					}
// 				}
// 			} else if reflect.ValueOf(val).Kind() == reflect.String && val != dataCheck.Get(key, "").(string) {
// 				resBool = false
// 			}
// 		}
// 	}

// 	return resBool
// }
