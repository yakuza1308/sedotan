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

	if g.serviceRunningStat == true {
		return errors.New("Service Already Running")
	}

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
