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
	NextGrabExe  time.Time
	LastGrabStat bool

	ServiceRunningStat bool

	ErrorNotes string

	//History||Summary
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
						}
						xN++
					}
					g.RowGrabbed += xN
					q.Close()
					g.DestDbox[key].IConnection.Close()

					g.Log.AddLog(fmt.Sprintf("[%s-%s] Fetch Data to destination finished, %d record fetch", g.Name, key, xN), "INFO")
				}
			}
		}
	}(g)
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
		g.Log.AddLog(fmt.Sprintf("[%s] Running Service", g.Name), "INFO")
		g.execService()
	} else {
		g.ErrorNotes = fmt.Sprintf("[%s] Running Service, Found : %s", g.Name, e)
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
		g.Log.AddLog(fmt.Sprintf("[%s] Stop Service", g.Name), "INFO")
	} else {
		g.Log.AddLog(fmt.Sprintf("[%s] Stop Service, Found : Service Not Running", g.Name), "ERROR")
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
