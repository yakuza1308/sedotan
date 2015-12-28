package modules

import (
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	"github.com/eaciit/knot/knot.v1"
	sdt "github.com/eaciit/sedotan/sedotan.v1"
	"github.com/eaciit/toolkit"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type GrabModule struct {
	nameId string
}

type StatService struct {
	name   string
	status bool
}

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/../"
	}()

	filename = wd + "data\\config.json"
	grabs    *sdt.GrabService
	grabber  *sdt.Grabber
)

// func NewGrabModule() *GrabModule {
// 	g := new(GrabModule)
// 	return g
// }

func Process(datas []interface{}) (string, bool) {
	var (
		status    string
		isServRun bool
	)
	for _, v := range datas {
		vToMap, _ := toolkit.ToM(v)
		if strings.ToLower(vToMap["calltype"].(string)) == "post" {
			grabs, _ = GrabPostConfig(vToMap)
		} else {
			grabs, _ = GrabGetConfig(vToMap)
		}
	}

	status, isServRun = StartService(grabs)
	return status, isServRun
}

func GrabPostConfig(data toolkit.M) (*sdt.GrabService, string) {
	var e error
	xGrabService := sdt.NewGrabService()
	xGrabService.Name = data["nameid"].(string) //"irondcecom"
	xGrabService.Url = data["url"].(string)     //"http://www.dce.com.cn/PublicWeb/MainServlet"

	xGrabService.SourceType = sdt.SourceType_Http

	xGrabService.GrabInterval = 5 * time.Minute
	xGrabService.TimeOutInterval = 1 * time.Minute //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, data["intervaltype"] /*"seconds"*/)

	grabConfig := sdt.Config{}

	dataurl := toolkit.M{}
	for _, grabconf := range data["grabconf"].(map[string]interface{}) {
		grabDataConf, _ := toolkit.ToM(grabconf)
		for key, subGrabDataConf := range grabDataConf {
			if reflect.ValueOf(subGrabDataConf).Kind() == reflect.Float64 {
				i := toolkit.ToInt(subGrabDataConf)
				toString := strconv.Itoa(i)
				dataurl[key] = toString
			} else {
				dataurl[key] = subGrabDataConf
			}
		}
	}
	// dataurl["Pu00231_Input.trade_date"] = "20151214"
	// dataurl["Pu00231_Input.variety"] = "i"
	// dataurl["Pu00231_Input.trade_type"] = "0"
	// dataurl["Submit"] = "Go"
	// dataurl["action"] = "Pu00231_result" //fmt.Println(i)

	grabConfig.SetFormValues(dataurl)

	xGrabService.ServGrabber = sdt.NewGrabber(xGrabService.Url, "POST", &grabConfig)

	logconfToMap, _ := toolkit.ToM(data["logconf"])
	logpath := logconfToMap["logpath"].(string)         //"E:\\data\\vale\\log"
	filename := logconfToMap["filename"].(string)       //"LOG-GRABDCETEST"
	filepattern := logconfToMap["filepattern"].(string) //"20060102"

	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
	if e != nil {
		return nil, e.Error()
	}

	xGrabService.Log = logconf

	xGrabService.ServGrabber.DataSettings = make(map[string]*sdt.DataSetting)
	xGrabService.DestDbox = make(map[string]*sdt.DestInfo)

	tempDataSetting := sdt.DataSetting{}
	tempDestInfo := sdt.DestInfo{}
	orCondition := []interface{}{}

	for _, dataSet := range data["datasettings"].([]interface{}) {
		dataToMap, _ := toolkit.ToM(dataSet)
		tempDataSetting.RowSelector = dataToMap["rowselector"].(string)

		for _, columnSet := range dataToMap["columnsettings"].([]interface{}) {
			columnToMap, _ := toolkit.ToM(columnSet)
			i := toolkit.ToInt(columnToMap["index"])
			tempDataSetting.Column(i, &sdt.GrabColumn{Alias: columnToMap["alias"].(string), Selector: columnToMap["selector"].(string)})
		}

		if hasRowdeletecond := dataToMap.Has("rowdeletecond"); hasRowdeletecond {
			for key, rowDeleteMap := range dataToMap["rowdeletecond"].(map[string]interface{}) {
				if key == "$or" || key == "$and" {
					for _, subDataRowDelete := range rowDeleteMap.([]interface{}) {
						for subIndex, getValueRowDelete := range subDataRowDelete.(map[string]interface{}) {
							tempDataSetting.Column(0, &sdt.GrabColumn{Alias: subIndex, Selector: getValueRowDelete.(string)})
						}
					}
					tempDataSetting.RowDeleteCond = toolkit.M{}.Set(key, orCondition)
				}
			}
		}
		xGrabService.ServGrabber.DataSettings[dataToMap["name"].(string)] = &tempDataSetting //DATA01 use name in datasettings

		connToMap, _ := toolkit.ToM(dataToMap["connectioninfo"])
		var db, usr, pwd string

		if hasDb := connToMap.Has("database"); !hasDb {
			db = ""
		} else {
			db = connToMap["database"].(string)
		}

		if hasUser := connToMap.Has("username"); !hasUser {
			usr = ""
		} else {
			usr = connToMap["username"].(string)
		}

		if hasPwd := connToMap.Has("password"); !hasPwd {
			pwd = ""
		} else {
			pwd = connToMap["username"].(string)
		}
		ci := dbox.ConnectionInfo{}
		ci.Host = connToMap["host"].(string) //"E:\\data\\vale\\Data_GrabIronTest.csv"
		ci.Database = db
		ci.UserName = usr
		ci.Password = pwd

		if hasSettings := connToMap.Has("settings"); !hasSettings {
			ci.Settings = nil
		} else {
			settingToMap, _ := toolkit.ToM(connToMap["settings"])
			ci.Settings = toolkit.M{}.Set("useheader", settingToMap["useheader"].(bool)).Set("delimiter", settingToMap["delimiter"].(string))
		}
		// ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")

		if hasCollection := connToMap.Has("collection"); !hasCollection {
			tempDestInfo.Collection = ""
		} else {
			tempDestInfo.Collection = connToMap["collection"].(string)
		}
		// tempDestInfo.Collection = ""
		tempDestInfo.Desttype = dataToMap["desttype"].(string)

		tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
		if e != nil {
			return nil, e.Error()
		}

		xGrabService.DestDbox[dataToMap["name"].(string)] = &tempDestInfo
	}

	return xGrabService, ""
}

func GrabGetConfig(data toolkit.M) (*sdt.GrabService, string) {
	xGrabService := sdt.NewGrabService()

	var (
		e error
		// isServiceRun bool
	)

	// isServiceRun = xGrabService.ServiceRunningStat
	// if isServiceRun {
	// 	return nil, "Service Running"
	// }

	xGrabService.Name = data["nameid"].(string) //"goldshfecom"
	xGrabService.Url = data["url"].(string)     //"http://www.shfe.com.cn/en/products/Gold/"

	xGrabService.SourceType = sdt.SourceType_Http

	xGrabService.GrabInterval = 5 * time.Minute
	xGrabService.TimeOutInterval = 1 * time.Minute //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, data["intervaltype"] /*"seconds"*/)

	grabConfig := sdt.Config{}

	xGrabService.ServGrabber = sdt.NewGrabber(xGrabService.Url, "GET", &grabConfig)

	logconfToMap, _ := toolkit.ToM(data["logconf"])
	logpath := logconfToMap["logpath"].(string)         //"E:\\data\\vale\\log"
	filename := logconfToMap["filename"].(string)       //"LOG-GRABSHFETEST"
	filepattern := logconfToMap["filepattern"].(string) //"20060102"

	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
	if e != nil {
		return nil, e.Error()
	}

	xGrabService.Log = logconf
	//==Data Setting and Destination Save =====================

	xGrabService.ServGrabber.DataSettings = make(map[string]*sdt.DataSetting)
	xGrabService.DestDbox = make(map[string]*sdt.DestInfo)

	// ==For Every Data Setting ===============================
	tempDataSetting := sdt.DataSetting{}
	tempDestInfo := sdt.DestInfo{}

	for _, dataSet := range data["datasettings"].([]interface{}) {
		dataToMap, _ := toolkit.ToM(dataSet)
		tempDataSetting.RowSelector = dataToMap["rowselector"].(string) //"#tab_conbox li:nth-child(1) .sjtable .listshuju tbody tr"

		for _, columnSet := range dataToMap["columnsettings"].([]interface{}) {
			columnToMap, _ := toolkit.ToM(columnSet)
			i := toolkit.ToInt(columnToMap["index"])
			tempDataSetting.Column(i, &sdt.GrabColumn{Alias: columnToMap["alias"].(string), Selector: columnToMap["selector"].(string)})
		}

		xGrabService.ServGrabber.DataSettings[dataToMap["name"].(string)] = &tempDataSetting //DATA01 use name in datasettings

		connToMap, _ := toolkit.ToM(dataToMap["connectioninfo"])
		var db, usr, pwd string

		if hasDb := connToMap.Has("database"); !hasDb {
			db = ""
		} else {
			db = connToMap["database"].(string)
		}

		if hasUser := connToMap.Has("username"); !hasUser {
			usr = ""
		} else {
			usr = connToMap["username"].(string)
		}

		if hasPwd := connToMap.Has("password"); !hasPwd {
			pwd = ""
		} else {
			pwd = connToMap["username"].(string)
		}

		ci := dbox.ConnectionInfo{}
		ci.Host = connToMap["host"].(string) //"E:\\WORKS\\data_test\\vale\\Data_Grab.csv"
		ci.Database = db
		ci.UserName = usr
		ci.Password = pwd

		if hasSettings := connToMap.Has("settings"); !hasSettings {
			ci.Settings = nil
		} else {
			settingToMap, _ := toolkit.ToM(connToMap["settings"])
			ci.Settings = toolkit.M{}.Set("useheader", settingToMap["useheader"].(bool)).Set("delimiter", settingToMap["delimiter"].(string))
		}

		if hasCollection := connToMap.Has("collection"); !hasCollection {
			tempDestInfo.Collection = ""
		} else {
			tempDestInfo.Collection = connToMap["collection"].(string)
		}

		tempDestInfo.Desttype = dataToMap["desttype"].(string)

		tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
		if e != nil {
			return nil, e.Error()
		}

		xGrabService.DestDbox[dataToMap["name"].(string)] = &tempDestInfo
	}

	return xGrabService, ""
}

func StartService(grabConfig *sdt.GrabService) (string, bool) {
	e := grabConfig.StartService()
	if e != nil {
		return e.Error(), grabConfig.ServiceRunningStat
	} else {
		if grabConfig.ServiceRunningStat {
			s := Insertstatservice(grabConfig.ServiceRunningStat, grabConfig.Name)

			ks := new(knot.Server)
			ks.Log().Info(fmt.Sprintf("==Start '%s' grab service==", s.(*StatService).name))
			knot.SharedObject().Set(s.(*StatService).name, grabConfig)
			return s.(*StatService).name, s.(*StatService).status
		}

		err, _ := StopService(grabConfig)
		if err != "" {
			return err, grabConfig.ServiceRunningStat
		}
	}

	return grabConfig.Name, grabConfig.ServiceRunningStat

}

func Insertstatservice(isServiceRun bool, nameid string) interface{} {
	var addserv interface{}
	if addserv == nil {
		addserv = &StatService{name: nameid, status: isServiceRun}
	} else {

	}
	return addserv
}

func StopProcess(datas []interface{}) (string, bool) {
	var (
		name       string
		isStopServ bool
	)
	for _, v := range datas {
		vToMap, _ := toolkit.ToM(v)
		if strings.ToLower(vToMap["calltype"].(string)) == "post" {

		} else {
			_, _ = GrabGetConfig(vToMap)
			name, isStopServ = StopService(knot.SharedObject().Get(vToMap["nameid"].(string)).(*sdt.GrabService))
		}
	}

	return name, isStopServ
}

func StopService(grabConfig *sdt.GrabService) (string, bool) {
	e := grabConfig.StopService()
	if e != nil {
		return e.Error(), grabConfig.ServiceRunningStat
	}
	ks := new(knot.Server)
	ks.Log().Info(fmt.Sprintf("==Stop '%s' grab service==", grabConfig.Name))
	return grabConfig.Name, grabConfig.ServiceRunningStat
}

func CheckStat(datas []interface{}) interface{} {
	var (
		grabStatus interface{}
	)
	for _, v := range datas {
		vToMap, _ := toolkit.ToM(v)
		if knot.SharedObject().Get(vToMap["nameid"].(string)) != nil {
			i := knot.SharedObject().Get(vToMap["nameid"].(string)).(*sdt.GrabService)

			grabStatus = ReadLog(vToMap["logconf"], i.ServiceRunningStat, i.Name)
		}
	}

	return grabStatus
}

func ReadLog(logConf interface{}, isRun bool, name string) interface{} {
	var grabsStatus = map[string]interface{}{}
	var grabStat bool

	dateNow := time.Now()
	logFile := fmt.Sprintf("%s\\%s%s!(EXTRA string=%s)", logConf.(map[string]interface{})["logpath"].(string), logConf.(map[string]interface{})["filename"].(string), "%", dateNow.Format(logConf.(map[string]interface{})["filepattern"].(string)))
	openLogFile, e := os.Open(logFile)
	if e != nil {
		grabsStatus["error"] = e.Error()
		return grabsStatus["error"]
	}
	defer openLogFile.Close()

	buf := make([]byte, 85)
	stat, e := os.Stat(logFile)
	start := stat.Size() - 85
	_, e = openLogFile.ReadAt(buf, start)
	if e != nil {
		grabsStatus["error"] = e.Error()
		return grabsStatus["error"]
	}

	splits := strings.Split(string(buf), "\n")
	var split []string
	if splits[0] != "" {
		split = strings.Split(splits[0], " ")
	} else {
		split = strings.Split(splits[1], " ")
	}

	if split[0] == "INFO" {
		grabStat = true
	} else {
		grabStat = false
	}

	grabsStatus["note"] = strings.Join(split[4:], " ")
	grabsStatus["dateNow"] = split[1]
	grabsStatus["timeNow"] = split[2]
	grabsStatus["status"] = grabStat
	grabsStatus["name"] = name
	grabsStatus["isRun"] = isRun

	return grabsStatus
}
