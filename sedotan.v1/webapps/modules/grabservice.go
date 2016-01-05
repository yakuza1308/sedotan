package modules

import (
	// "bufio"
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
	// "strings"
	"time"
)

type GrabModule struct {
	nameId, lastGrabTime string
	gi, ti               time.Time
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

func NewGrabService() *GrabModule {
	g := new(GrabModule)
	return g
}

func Process(datas []interface{}) (string, bool) {
	var (
		status    string
		isServRun bool
	)
	for _, v := range datas {
		vToMap, _ := toolkit.ToM(v)
		grabs, _ = GrabConfig(vToMap)
	}

	status, isServRun = StartService(grabs)

	return status, isServRun
}

func GrabConfig(data toolkit.M) (*sdt.GrabService, string) {
	var e error
	var gi, ti time.Duration
	xGrabService := sdt.NewGrabService()
	xGrabService.Name = data["nameid"].(string) //"irondcecom"
	xGrabService.Url = data["url"].(string)     //"http://www.dce.com.cn/PublicWeb/MainServlet"

	xGrabService.SourceType = sdt.SourceType_Http

	grabintervalToInt := toolkit.ToInt(data["grabinterval"])
	timeintervalToInt := toolkit.ToInt(data["timeoutinterval"])
	if data["intervaltype"].(string) == "seconds" {
		gi = time.Duration(grabintervalToInt) * time.Second
		ti = time.Duration(timeintervalToInt) * time.Second
	} else if data["intervaltype"].(string) == "minutes" {
		gi = time.Duration(grabintervalToInt) * time.Minute
		ti = time.Duration(timeintervalToInt) * time.Minute
	} else if data["intervaltype"].(string) == "hours" {
		gi = time.Duration(grabintervalToInt) * time.Hour
		ti = time.Duration(timeintervalToInt) * time.Hour
	}

	xGrabService.GrabInterval = gi    //* time.Minute
	xGrabService.TimeOutInterval = ti //* time.Minute //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", timeintervalToInt, data["intervaltype"] /*"seconds"*/)

	grabConfig := sdt.Config{}

	if data["calltype"].(string) == "POST" {
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

		grabConfig.SetFormValues(dataurl)
	}

	xGrabService.ServGrabber = sdt.NewGrabber(xGrabService.Url, data["calltype"].(string), &grabConfig)

	logconfToMap, _ := toolkit.ToM(data["logconf"])
	logpath := logconfToMap["logpath"].(string)           //"E:\\data\\vale\\log"
	filename := logconfToMap["filename"].(string) + "-%s" //"LOG-GRABDCETEST"
	filepattern := logconfToMap["filepattern"].(string)   //"20060102"

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
	var condition string

	for _, dataSet := range data["datasettings"].([]interface{}) {
		dataToMap, _ := toolkit.ToM(dataSet)
		tempDataSetting.RowSelector = dataToMap["rowselector"].(string)

		for _, columnSet := range dataToMap["columnsettings"].([]interface{}) {
			columnToMap, _ := toolkit.ToM(columnSet)
			i := toolkit.ToInt(columnToMap["index"])
			tempDataSetting.Column(i, &sdt.GrabColumn{Alias: columnToMap["alias"].(string), Selector: columnToMap["selector"].(string)})
		}

		if data["calltype"].(string) == "POST" {
			// orCondition := []interface{}{}
			if hasRowdeletecond := dataToMap.Has("rowdeletecond"); hasRowdeletecond {
				for key, rowDeleteMap := range dataToMap["rowdeletecond"].(map[string]interface{}) {
					if key == "$or" || key == "$and" {
						for _, subDataRowDelete := range rowDeleteMap.([]interface{}) {
							for subIndex, getValueRowDelete := range subDataRowDelete.(map[string]interface{}) {
								orCondition = append(orCondition, map[string]interface{}{subIndex: getValueRowDelete})
							}
						}
						condition = key
					}
				}
			}
			tempDataSetting.RowDeleteCond = toolkit.M{}.Set(condition, orCondition)
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
			ci.Settings = toolkit.M{}.Set("useheader", settingToMap["useheader"].(bool)).Set("delimiter", settingToMap["delimiter"])
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
		name, isStopServ = StopService(knot.SharedObject().Get(vToMap["nameid"].(string)).(*sdt.GrabService))
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

func (g *GrabModule) CheckStat(datas []interface{}) interface{} {
	var (
		grabStatus         interface{}
		lastDate, nextDate string
	)

	for _, v := range datas {
		vToMap, _ := toolkit.ToM(v)

		if knot.SharedObject().Get(vToMap["nameid"].(string)) != nil {
			i := knot.SharedObject().Get(vToMap["nameid"].(string)).(*sdt.GrabService)
			tLast := i.LastGrabExe.Format("2006/01/02 15:04:05")
			if tLast != "0001/01/01 00:00:00" {
				lastDate = tLast
			}
			tNext := i.NextGrabExe.Format("2006/01/02 15:04:05")
			if tNext != "0001/01/01 00:00:00" {
				nextDate = tNext
			}
			grabStatus = ReadLog(vToMap["logconf"], i.ServiceRunningStat, i.Name, lastDate, nextDate, i.LastGrabStat, i.ErrorNotes)
		} else {
			grabStatus = ReadLog(vToMap["logconf"], false, vToMap["nameid"].(string), "", "", false, "")
		}
	}

	return grabStatus
}

func ReadLog(logConf interface{}, isRun bool, name string, lastDate string, nextDate string, lastGrab bool, errorNote string) interface{} {
	var grabsStatus = map[string]interface{}{}
	// var grabStat bool

	grabsStatus["errorNote"] = errorNote
	grabsStatus["lastDate"] = lastDate
	grabsStatus["nextDate"] = nextDate
	grabsStatus["grabStat"] = lastGrab
	grabsStatus["name"] = name
	grabsStatus["isRun"] = isRun
	// dateNow := time.Now()
	// logFile := fmt.Sprintf("%s\\%s-%s", logConf.(map[string]interface{})["logpath"].(string), logConf.(map[string]interface{})["filename"].(string), dateNow.Format(logConf.(map[string]interface{})["filepattern"].(string)))
	// // logFile := fmt.Sprintf("%s\\%s", logConf.(map[string]interface{})["logpath"].(string), logConf.(map[string]interface{})["filename"].(string))

	// openLogFile, e := os.Open(logFile)
	// if e != nil {
	// 	if b := strings.Contains(e.Error(), "The system cannot find the file specified"); b {
	// 		grabsStatus["grabStat"] = b
	// 	} else {
	// 		grabsStatus["grabStat"] = false
	// 	}
	// 	grabsStatus["error"] = e.Error()
	// 	grabsStatus["name"] = name

	// 	fmt.Printf("grabsStatus error:%v\n", grabsStatus)
	// 	return grabsStatus
	// }
	// defer openLogFile.Close()

	// scanner := bufio.NewScanner(openLogFile)
	// i := 0
	// var lastLine string
	// for scanner.Scan() {
	// 	i++
	// 	lastLine = scanner.Text()

	// 	lastNum := i - 1
	// 	if i == lastNum {
	// 		break
	// 	}
	// }

	// if lastLine != "" {
	// 	splits := strings.Split(lastLine, " ")
	// 	if splits[0] == "INFO" {
	// 		grabStat = true
	// 	} else if splits[0] == "ERROR" {
	// 		grabStat = false
	// 	}

	// 	grabsStatus["note"] = strings.Join(splits[4:], " ")
	// 	grabsStatus["dateTimeNow"] = lastDate
	// 	// grabsStatus["timeNow"] = splits[2]
	// 	grabsStatus["grabStat"] = grabStat
	// }
	// grabsStatus["name"] = name
	// grabsStatus["isRun"] = isRun
	// fmt.Printf("grabsStatus:%v\n", grabsStatus)
	return grabsStatus
}

// func OpenLog(dateString string) string {
// 	// buka file
// 	var file, err = os.OpenFile(path, os.O_RDONLY, 0644)
// 	checkError(err)
// 	defer file.Close()

// 	// baca file
// 	var text = make([]byte, 1024)
// 	for {
// 		n, err := file.Read(text)
// 		if err != io.EOF {
// 			checkError(err)
// 		}
// 		if n == 0 {
// 			break
// 		}
// 	}
// 	fmt.Println(string(text))
// 	checkError(err)
// }
