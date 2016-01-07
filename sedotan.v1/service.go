package sedotan

import (
	"errors"
	"fmt"
	"github.com/eaciit/cast"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"os"
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
	HistoryPath         string
	HistoryRecPath      string
	Log                 *toolkit.LogEngine

	ServGrabber *Grabber

	LastGrabExe  time.Time
	NextGrabExe  time.Time
	LastGrabStat bool

	ServiceRunningStat bool

	ErrorNotes string

	//Snapshot
	StartDate  time.Time
	EndDate    time.Time
	GrabCount  int
	RowGrabbed int
	ErrorFound int
}

type DestInfo struct {
	dbox.IConnection
	Collection string
	Desttype   string
}

func NewGrabService() *GrabService {
	g := new(GrabService)
	g.SourceType = SourceType_Http
	g.GrabInterval = 5 * time.Minute
	g.TimeOutInterval = 1 * time.Minute
	g.ServiceRunningStat = false

	dir, _ := os.Getwd()
	g.HistoryPath = strings.Replace(dir, " ", "\\ ", -1)
	g.HistoryRecPath = strings.Replace(dir, " ", "\\ ", -1)
	return g
}

func (g *GrabService) execService() {
	g.LastGrabStat = false
	go func(g *GrabService) {
		for g.ServiceRunningStat {

			if g.LastGrabStat {
				<-time.After(g.GrabInterval)
			} else {
				<-time.After(g.TimeOutInterval)
			}

			if !g.ServiceRunningStat {
				continue
			}

			g.ErrorNotes = ""
			g.LastGrabExe = time.Now()
			g.NextGrabExe = time.Now().Add(g.GrabInterval)
			g.LastGrabStat = true
			g.Log.AddLog(fmt.Sprintf("[%s] Grab Started %s", g.Name, g.Url), "INFO")
			g.GrabCount += 1

			if e := g.ServGrabber.Grab(nil); e != nil {
				g.ErrorNotes = fmt.Sprintf("[%s] Grab Failed %s, repeat after %s :%s", g.Name, g.Url, g.TimeOutIntervalInfo, e)
				g.Log.AddLog(g.ErrorNotes, "ERROR")
				g.NextGrabExe = time.Now().Add(g.TimeOutInterval)
				g.LastGrabStat = false
				g.ErrorFound += 1
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
						g.ErrorNotes = fmt.Sprintf("[%s-%s] Fetch Result Failed : ", g.Name, key, e)
						g.Log.AddLog(g.ErrorNotes, "ERROR")
						continue
					}

					e = g.DestDbox[key].IConnection.Connect()
					if e != nil {
						g.ErrorNotes = fmt.Sprintf("[%s-%s] Connect to destination failed [%s-%s]:%s", g.Name, key, g.DestDbox[key].Desttype, g.DestDbox[key].IConnection.Info().Host, e)
						g.Log.AddLog(g.ErrorNotes, "ERROR")
						continue
					}

					var q dbox.IQuery
					if g.DestDbox[key].Collection == "" {
						q = g.DestDbox[key].IConnection.NewQuery().SetConfig("multiexec", true).Save()
					} else {
						q = g.DestDbox[key].IConnection.NewQuery().SetConfig("multiexec", true).From(g.DestDbox[key].Collection).Save()
					}
					xN := 0
					iN := 0
					for _, doc := range docs {
						for key, val := range doc {
							doc[key] = strings.TrimSpace(fmt.Sprintf("%s", val))
						}

						if g.DestDbox[key].Desttype == "mongo" {
							doc["_id"] = toolkit.GenerateRandomString("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnpqrstuvwxyz", 32)
						}

						e = q.Exec(toolkit.M{
							"data": doc,
						})

						if e != nil {
							g.ErrorNotes = fmt.Sprintf("[%s-%s] Unable to insert [%s-%s]:%s", g.Name, key, g.DestDbox[key].Desttype, g.DestDbox[key].IConnection.Info().Host, e)
							g.Log.AddLog(g.ErrorNotes, "ERROR")
							g.ErrorFound += 1
						} else {
							iN += 1
						}
						xN++
					}
					g.RowGrabbed += xN
					q.Close()
					g.DestDbox[key].IConnection.Close()

					g.Log.AddLog(fmt.Sprintf("[%s-%s] Fetch Data to destination finished with %d record fetch", g.Name, key, xN), "INFO")

					//ADD History
					recfile := g.AddRecHistory(docs)
					historyservice := toolkit.M{}.Set("datasettingname", key).Set("grabdate", g.LastGrabExe).Set("rowgrabbed", g.RowGrabbed).
						Set("rowsaved", iN).Set("note", g.ErrorNotes).Set("grabstatus", "FAILED").Set("recfile", recfile)
					if g.LastGrabStat {
						historyservice.Set("grabstatus", "SUCCESS")
					}

					g.AddHistory(historyservice)
				}
			}
		}
	}(g)
}

func (g *GrabService) AddRecHistory(docs []toolkit.M) string {
	var config = map[string]interface{}{"useheader": true, "delimiter": ",", "newfile": true}
	file := fmt.Sprintf("%s%s-%s.csv", g.HistoryRecPath, g.Name, cast.Date2String(time.Now(), "YYYYMMddHHmmss"))
	ci := &dbox.ConnectionInfo{file, "", "", "", config}
	c, e := dbox.NewConnection("csv", ci)
	if e != nil {
		g.ErrorNotes = fmt.Sprintf("[%s] Setup connection to Record history failed [csv-%s]:%s", g.Name, file, e)
		g.Log.AddLog(g.ErrorNotes, "ERROR")
		return ""
	}

	e = c.Connect()
	if e != nil {
		g.ErrorNotes = fmt.Sprintf("[%s] Setup connection to history failed [csv-%s]:%s", g.Name, file, e)
		g.Log.AddLog(g.ErrorNotes, "ERROR")
		return ""
	}

	// q := c.NewQuery().SetConfig("multiexec", true).Save()

	for _, doc := range docs {
		e = c.NewQuery().Insert().Exec(toolkit.M{"data": doc})
		if e != nil {
			g.ErrorNotes = fmt.Sprintf("[%s] Insert to history failed [csv-%s]:%s", g.Name, file, e)
			g.Log.AddLog(g.ErrorNotes, "ERROR")
			return ""
		}
	}
	c.Close()

	return file
}

func (g *GrabService) AddHistory(history toolkit.M) {
	mapHeader := make([]toolkit.M, 7)
	mapHeader[0] = toolkit.M{}.Set("datasettingname", "string")
	mapHeader[1] = toolkit.M{}.Set("grabdate", "date")
	mapHeader[2] = toolkit.M{}.Set("grabstatus", "string")
	mapHeader[3] = toolkit.M{}.Set("rowgrabbed", "int")
	mapHeader[4] = toolkit.M{}.Set("rowsaved", "int")
	mapHeader[5] = toolkit.M{}.Set("note", "string")
	mapHeader[6] = toolkit.M{}.Set("recfile", "string")

	var config = map[string]interface{}{"mapheader": mapHeader, "useheader": true, "delimiter": ",", "newfile": true}
	file := fmt.Sprintf("%s%s-%s.csv", g.HistoryPath, g.Name, cast.Date2String(time.Now(), "YYYYMM"))
	ci := &dbox.ConnectionInfo{file, "", "", "", config}
	c, e := dbox.NewConnection("csv", ci)
	if e != nil {
		g.ErrorNotes = fmt.Sprintf("[%s] Setup connection to history failed [csv-%s]:%s", g.Name, file, e)
		g.Log.AddLog(g.ErrorNotes, "ERROR")
		return
	}

	e = c.Connect()
	if e != nil {
		g.ErrorNotes = fmt.Sprintf("[%s] Setup connection to history failed [csv-%s]:%s", g.Name, file, e)
		g.Log.AddLog(g.ErrorNotes, "ERROR")
		return
	}

	e = c.NewQuery().Insert().Exec(toolkit.M{"data": history})
	if e != nil {
		g.ErrorNotes = fmt.Sprintf("[%s] Insert to history failed [csv-%s]:%s", g.Name, file, e)
		g.Log.AddLog(g.ErrorNotes, "ERROR")
		c.Close()
		return
	}

	c.Close()
}

func (g *GrabService) StartService() error {
	g.ErrorNotes = ""
	if g.ServiceRunningStat == true {
		return errors.New("Service Already Running")
	}

	g.ServiceRunningStat = false
	noErrorFound, e := g.validateService()

	if noErrorFound {
		g.StartDate = time.Now()
		g.EndDate = time.Time{}
		g.GrabCount = 0
		g.RowGrabbed = 0
		g.ErrorFound = 0

		g.ServiceRunningStat = true
		g.Log.AddLog(fmt.Sprintf("[%s] Service Running", g.Name), "INFO")
		g.execService()
	} else {
		g.ErrorNotes = fmt.Sprintf("[%s] Service Running, Found : %s", g.Name, e)
		g.Log.AddLog(g.ErrorNotes, "ERROR")
		g.ErrorFound += 1
		return e
	}

	return nil
}

func (g *GrabService) StopService() error {
	if g.ServiceRunningStat {
		g.EndDate = time.Now()
		g.ServiceRunningStat = false
		g.Log.AddLog(fmt.Sprintf("[%s] Service Stop", g.Name), "INFO")
	} else {
		g.Log.AddLog(fmt.Sprintf("[%s] Service Stop, Found : Service Not Running", g.Name), "ERROR")
		g.ErrorFound += 1
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
