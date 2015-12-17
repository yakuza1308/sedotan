package sedotan

import (
	"errors"
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"strings"
	"time"
)

type SourceTypeEnum int

const (
	SourceType_Http SourceTypeEnum = iota
)

type GrabService struct {
	Name                string
	Url                 string
	SourceType          SourceTypeEnum
	GrabInterval        time.Duration
	TimeOutInterval     time.Duration
	TimeOutIntervalInfo string
	DestDbox            map[string]*DestInfo
	Log                 *toolkit.LogEngine

	ServGrabber *Grabber

	LastGrabExe  time.Time
	LastGrabStat bool

	serviceRunningStat bool
}

type DestInfo struct {
	dbox.IConnection
	collection string
	desttype   string
}

func newGrabService() *GrabService {
	g := new(GrabService)
	g.SourceType = SourceType_Http
	g.GrabInterval = 5 * time.Minute
	g.TimeOutInterval = 1 * time.Minute
	return g
}

func (g *GrabService) execService() {
	g.LastGrabStat = false
	go func(g *GrabService) {
		for g.serviceRunningStat {

			if g.LastGrabStat {
				<-time.After(g.GrabInterval)
			} else {
				<-time.After(g.TimeOutInterval)
			}

			g.LastGrabExe = time.Now()
			g.LastGrabStat = true
			g.Log.AddLog(fmt.Sprintf("[%s] Grab Started %s", g.Name, g.Url), "INFO")

			if e := g.ServGrabber.Grab(nil); e != nil {
				g.Log.AddLog(fmt.Sprintf("[%s] Grab Failed %s, repeat after %s :%s", g.Name, g.Url, g.TimeOutIntervalInfo, e), "ERROR")
				g.LastGrabStat = false
				continue
			} else {
				g.Log.AddLog(fmt.Sprintf("[%s] Grab Success %s", g.Name, g.Url), "INFO")
			}

			if g.LastGrabStat {
				for key, _ := range g.ServGrabber.Config.DataSettings {
					g.Log.AddLog(fmt.Sprintf("[%s-%s] Fetch Data to destination started", g.Name, key), "INFO")

					docs := []toolkit.M{}
					e := g.ServGrabber.ResultFromHtml(key, &docs)
					if e != nil {
						g.Log.AddLog(fmt.Sprintf("[%s-%s] Fetch Result Failed : ", g.Name, key, e), "ERROR")
						continue
					}

					e = g.DestDbox[key].IConnection.Connect()
					if e != nil {
						g.Log.AddLog(fmt.Sprintf("[%s-%s] Connect to destination failed [%s-%s]:%s", g.Name, key, g.DestDbox[key].desttype, g.DestDbox[key].IConnection.Info().Host, e), "ERROR")
						continue
					}

					var q dbox.IQuery
					if g.DestDbox[key].collection == "" {
						q = g.DestDbox[key].IConnection.NewQuery().SetConfig("multiexec", true).Save()
					} else {
						q = g.DestDbox[key].IConnection.NewQuery().SetConfig("multiexec", true).From(g.DestDbox[key].collection).Save()
					}

					for x, doc := range docs {
						for key, val := range doc {
							doc[key] = strings.TrimSpace(fmt.Sprintf("%s", val))
						}

						if g.DestDbox[key].desttype == "mongo" {
							doc["_id"] = x
						}

						e = q.Exec(toolkit.M{
							"data": doc,
						})
						if e != nil {
							g.Log.AddLog(fmt.Sprintf("[%s-%s] Unable to insert [%s-%s]:%s", g.Name, key, g.DestDbox[key].desttype, g.DestDbox[key].IConnection.Info().Host, e), "ERROR")
						}
					}
					q.Close()
					g.DestDbox[key].IConnection.Close()

					g.Log.AddLog(fmt.Sprintf("[%s-%s] Fetch Data to destination finished", g.Name, key), "INFO")
				}
			}
		}
	}(g)
}

func (g *GrabService) StartService() error {
	g.serviceRunningStat = false
	noErrorFound, e := g.validateService()

	if noErrorFound {
		g.serviceRunningStat = true
		g.Log.AddLog(fmt.Sprintf("[%s] Running Service", g.Name), "INFO")
		g.execService()
	} else {
		return e
	}

	return nil
}

func (g *GrabService) StopService() error {
	if g.serviceRunningStat {
		g.serviceRunningStat = false
		g.Log.AddLog(fmt.Sprintf("[%s] Stop Service", g.Name), "INFO")
	} else {
		return errors.New("Service Not Running")
	}

	return nil
}

func (g *GrabService) validateService() (bool, error) {
	if g.Log == nil {
		return false, errors.New("Log Not Found")
	}

	if g.Name == "" {
		return false, errors.New("Name Not Found")
	}

	if g.SourceType != SourceType_Http {
		return false, errors.New("Source Type Not Set")
	}

	if g.Url == "" {
		return false, errors.New("Url Not Found")
	}

	for key, val := range g.DestDbox {
		e := val.IConnection.Connect()
		if e != nil {
			return false, errors.New(fmt.Sprintf("[%s] Found : %s", key, e))
		}
	}

	// Do Validate
	return true, nil
}

// func main() {
// 	jsonpath := "E:\\data\\vale\\config.json"

// 	ci := &dbox.ConnectionInfo{jsonpath, "", "", "", nil}
// 	c, e := dbox.NewConnection("json", ci)
// 	if e != nil {
// 		fmt.Println("Error Found : ", e)
// 	}

// 	e = c.Connect()
// 	if e != nil {
// 		fmt.Println("Error Found : ", e)
// 	}

// 	csr, e := c.NewQuery().Select("nameid", "url", "authtype").Cursor(nil)
// 	ds, e := csr.Fetch(nil, 0, false)

// 	for x, val := range ds.Data {
// 		mapVal, e := toolkit.ToM(val.(map[string]interface{}))
// 		errorNotFound := true

// 		logpath := ""
// 		filename := ""
// 		filepattern := ""
// 		logtofile := false
// 		logtostdout := true

// 		if mapVal.Has("logconf") {
// 			tempLogConf, e := toolkit.ToM(mapVal.Get("logconf", nil).(map[string]interface{}))
// 			if e != nil {
// 				fmt.Println("Error Read Log Conf : %s", e.Error())
// 			}

// 			logpath = tempLogConf.Get("logpath", "").(string)
// 			filename = tempLogConf.Get("filename", "").(string)
// 			filepattern = tempLogConf.Get("filepattern", "").(string)

// 			if logpath != "" && filename != "" {
// 				logtofile = true
// 				logtostdout = false
// 			}
// 		}

// 		logconf, e := toolkit.NewLog(logtostdout, logtofile, logpath, filename, filepattern)

// 		if e != nil {
// 			fmt.Println("Error Start Log : %s", e.Error())
// 		}

// 		logconf.AddLog(fmt.Sprintf("[READCONF-%d] Reading next configuration", x), "INFO")

// 		if e != nil {
// 			logconf.AddLog(fmt.Sprintf("[READCONF-%d] Found : %s", x, e), "ERROR")
// 			errorNotFound = false
// 			continue
// 		}

// 		if !(mapVal.Has("nameid")) {
// 			logconf.AddLog(fmt.Sprintf("[READCONF-%d] Name ID not found, skip to next config", x), "ERROR")
// 			errorNotFound = false
// 			continue
// 		}

// 		tempGrab := newGrabService()
// 		tempGrab.Log = logconf
// 		tempGrab.Name = mapVal.Get("nameid", "").(string)
// 		logconf.AddLog(fmt.Sprintf("[READCONF-%d] Loading Configuration - %s", x, tempGrab.Name), "INFO")
// 		if mapVal.Get("sourcetype", "").(string) == "SourceType_Http" {
// 			tempGrab.SourceType = SourceType_Http
// 		}

// 		tempGrab.Url = mapVal.Get("url", "").(string)
// 		if tempGrab.Url == "" {
// 			logconf.AddLog(fmt.Sprintf("[READCONF-%d] Found : %s", x, "URL Not Found"), "ERROR")
// 			errorNotFound = false
// 			continue
// 		}

// 		if !(mapVal.Has("calltype")) {
// 			logconf.AddLog(fmt.Sprintf("[READCONF-%d] Found : %s", x, "Call Type Not Found"), "ERROR")
// 			errorNotFound = false
// 			continue
// 		}

// 		if mapVal.Has("grabconfig") {
// 			mapValConfig, e := toolkit.ToM(mapVal.Get("grabconfig", nil).(map[string]interface{}))
// 			if e != nil {
// 				logconf.AddLog(fmt.Sprintf("[READCONF-%d] Grab Config Found : %s", x, e), "ERROR")
// 			}

// 			GrabConfig := sedotan.Config{}
// 			if mapValConfig.Has("formvalues") {
// 				tempFormValues, e := toolkit.ToM(mapValConfig.Get("formvalues", nil).(map[string]interface{}))
// 				if e != nil {
// 					logconf.AddLog(fmt.Sprintf("[READCONF-%d] Grab Config Form Values Found : %s", x, e), "ERROR")
// 				}
// 				GrabConfig.SetFormValues(tempFormValues)
// 			}

// 			GrabConfig.URL = mapValConfig.Get("url", "").(string)
// 			GrabConfig.CallType = mapValConfig.Get("calltype", "").(string)

// 			GrabConfig.AuthType = mapValConfig.Get("authtype", "").(string)
// 			GrabConfig.AuthUserId = mapValConfig.Get("authuserid", "").(string)
// 			GrabConfig.AuthPassword = mapValConfig.Get("authpassword", "").(string)

// 			tempGrab.ServGrabber = sedotan.NewGrabber(tempGrab.Url, mapVal.Get("calltype", "").(string), &GrabConfig)
// 		} else {
// 			tempGrab.ServGrabber = sedotan.NewGrabber(tempGrab.Url, mapVal.Get("calltype", "").(string), nil)
// 		}

// 		if mapVal.Has("intervaltype") {
// 			intervaltype := mapVal.Get("intervaltype", "").(string)
// 			grabinterval := mapVal.Get("grabinterval", 1).(float64)
// 			timeoutinterval := mapVal.Get("timeoutinterval", 1).(float64)

// 			tempGrab.TimeOutIntervalInfo = fmt.Sprintf("%v %s", timeoutinterval, intervaltype)
// 			switch intervaltype {
// 			case "hours":
// 				tempGrab.GrabInterval = time.Duration(grabinterval) * time.Hour
// 				tempGrab.TimeOutInterval = time.Duration(timeoutinterval) * time.Hour
// 			case "minutes":
// 				tempGrab.GrabInterval = time.Duration(grabinterval) * time.Minute
// 				tempGrab.TimeOutInterval = time.Duration(timeoutinterval) * time.Minute
// 			default:
// 				tempGrab.GrabInterval = time.Duration(grabinterval) * time.Second
// 				tempGrab.TimeOutInterval = time.Duration(timeoutinterval) * time.Second
// 				tempGrab.TimeOutIntervalInfo = fmt.Sprintf("%v %s", timeoutinterval, "seconds")
// 			}
// 		}

// 		if mapVal.Has("datasettings") {
// 			logconf.AddLog(fmt.Sprintf("[READCONF-%s] Loading Data Settings", tempGrab.Name), "INFO")

// 			tempGrab.ServGrabber.DataSettings = make(map[string]*sedotan.DataSetting)
// 			tempGrab.DestDbox = make(map[string]*DestInfo)

// 			for xX, xVal := range mapVal["datasettings"].([]interface{}) {
// 				logconf.AddLog(fmt.Sprintf("[READCONF-%s] Data Settings - %d", tempGrab.Name, xX), "INFO")
// 				tempDataSetting := sedotan.DataSetting{}

// 				mapxVal, e := toolkit.ToM(xVal.(map[string]interface{}))

// 				if e != nil {
// 					logconf.AddLog(fmt.Sprintf("[READCONF-%s] Found : %s (DS[%d])", tempGrab.Name, e, xX), "ERROR")
// 					errorNotFound = false
// 					continue
// 				}

// 				tempDataSetting.RowSelector = mapxVal.Get("rowselector", "").(string)
// 				for xXX, xxVal := range mapxVal["columnsettings"].([]interface{}) {
// 					mapxxVal, e := toolkit.ToM(xxVal.(map[string]interface{}))
// 					if e != nil {
// 						logconf.AddLog(fmt.Sprintf("[READCONF-%s] Found : %s (DS:ColS[%d:%d])", tempGrab.Name, e, xX, xXX), "ERROR")
// 						errorNotFound = false
// 						continue
// 					}

// 					tempGrabColumn := sedotan.GrabColumn{}
// 					tempGrabColumn.Alias = mapxxVal.Get("alias", "").(string)
// 					tempGrabColumn.Selector = mapxxVal.Get("selector", "").(string)

// 					if mapVal.Has("valuetype") {
// 						tempGrabColumn.Selector = mapxxVal.Get("valuetype", "").(string)
// 					}

// 					if mapVal.Has("attrname") {
// 						tempGrabColumn.Selector = mapxxVal.Get("attrname", "").(string)
// 					}

// 					tempDataSetting.Column(int(mapxxVal.Get("index", 0).(float64)), &tempGrabColumn)
// 				}

// 				if mapxVal.Has("rowdeletecond") {
// 					tempDataSetting.RowDeleteCond, e = toolkit.ToM(mapxVal.Get("rowdeletecond", nil).(map[string]interface{}))
// 					if e != nil {
// 						logconf.AddLog(fmt.Sprintf("[READCONF-%s] Delete Condition Found : %s (DS[%d])", tempGrab.Name, e, xX), "ERROR")
// 					}
// 				}

// 				tempGrab.ServGrabber.DataSettings[mapxVal.Get("name", "").(string)] = &tempDataSetting

// 				if !(mapxVal.Has("desttype")) {
// 					logconf.AddLog(fmt.Sprintf("[READCONF-%s] Destination Not Found : %s (DS[%d])", tempGrab.Name, e, xX), "ERROR")
// 					errorNotFound = false
// 					continue
// 				}

// 				if mapxVal.Has("connectioninfo") {
// 					ci := dbox.ConnectionInfo{}
// 					mapConnVal, e := toolkit.ToM(mapxVal.Get("connectioninfo", nil).(map[string]interface{}))
// 					if e != nil {
// 						logconf.AddLog(fmt.Sprintf("[READCONF-%s] Found : %s (DS[%d],DestConn)", tempGrab.Name, e, xX), "ERROR")
// 						errorNotFound = false
// 						continue
// 					}

// 					ci.Host = mapConnVal.Get("host", "").(string)
// 					ci.Database = mapConnVal.Get("database", "").(string)
// 					ci.UserName = mapConnVal.Get("userName", "").(string)
// 					ci.Password = mapConnVal.Get("password", "").(string)

// 					tempSetting := mapConnVal.Get("settings", nil)
// 					if tempSetting != nil {
// 						ci.Settings, e = toolkit.ToM(tempSetting.(map[string]interface{}))
// 						if e != nil {
// 							logconf.AddLog(fmt.Sprintf("[READCONF-%s] Found : %s (DS[%d],DestConnSettings)", tempGrab.Name, e, xX), "ERROR")
// 							errorNotFound = false
// 							continue
// 						}
// 					} else {
// 						ci.Settings = nil
// 					}

// 					tempDestInfo := DestInfo{}
// 					tempDestInfo.IConnection, e = dbox.NewConnection(mapxVal.Get("desttype", "").(string), &ci)
// 					if e != nil {
// 						logconf.AddLog(fmt.Sprintf("[READCONF-%s] Found : %s (DS[%d],DestNewConnCreate)", tempGrab.Name, e, xX), "ERROR")
// 						errorNotFound = false
// 						continue
// 					}
// 					tempDestInfo.collection = mapConnVal.Get("collection", "").(string)
// 					tempDestInfo.desttype = mapxVal.Get("desttype", "").(string)

// 					tempGrab.DestDbox[mapxVal.Get("name", "").(string)] = &tempDestInfo

// 					e = tempGrab.DestDbox[mapxVal.Get("name", "").(string)].IConnection.Connect()
// 					if e != nil {
// 						logconf.AddLog(fmt.Sprintf("[READCONF-%s] Found : %s (DS[%d],DestConnTest)", tempGrab.Name, e, xX), "ERROR")
// 						errorNotFound = false
// 						continue
// 					} else {
// 						logconf.AddLog(fmt.Sprintf("[READCONF-%s] Test Connect Success (DS[%d]) [%s-%s]", tempGrab.Name, xX, tempDestInfo.desttype, ci.Host), "INFO")
// 					}
// 					tempGrab.DestDbox[mapxVal.Get("name", "").(string)].IConnection.Close()
// 				}
// 			}
// 		}

// 		if errorNotFound {
// 			tempGrab.Exec()
// 		}
// 	}

// 	c.Close()

// 	for {
// 		fmt.Printf(".")
// 		time.Sleep(3000 * time.Millisecond)
// 	}
// }
