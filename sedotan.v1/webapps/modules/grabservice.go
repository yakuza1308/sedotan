package modules

import (
	"fmt"
	// "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	// sdt "github.com/eaciit/sedotan/sedotan.v1"
	// "github.com/eaciit/toolkit"
	// "time"
)

type GrabModule struct {
	nameId string
}

func NewGrabModule() *GrabModule {
	g := new(GrabModule)
	return g
}

func Grabget(datas []interface{}) string {
	for i, v := range datas {
		fmt.Printf("index:%v\n", i)
		fmt.Printf("datas:%v\n", v.(map[string]interface{})["nameid"])
	}

	return "hi"
	/*
		xGrabService := sdt.NewGrabService()
		xGrabService.Name = "goldshfecom"
		xGrabService.Url = "http://www.shfe.com.cn/en/products/Gold/"

		xGrabService.SourceType = sdt.SourceType_Http

		xGrabService.GrabInterval = 5 * time.Minute
		xGrabService.TimeOutInterval = 1 * time.Minute //time.Hour, time.Minute, time.Second

		xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, "seconds")

		grabConfig := sdt.Config{}
		// if has grabconfig

		// CallType     string
		// FormValues   toolkit.M
		// AuthType     string
		// AuthUserId   string
		// AuthPassword string
		//==================

		xGrabService.ServGrabber = sdt.NewGrabber(xGrabService.Url, "GET", &grabConfig)

		//===================================================================

		//==For Data Log          ===========================================

		// 	logpath = tempLogConf.Get("logpath", "").(string)
		// 	filename = tempLogConf.Get("filename", "").(string)
		// 	filepattern = tempLogConf.Get("filepattern", "").(string)

		logpath := "E:\\data\\vale\\log"
		filename := "LOG-GRABSHFETEST"
		filepattern := "20060102"

		logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
		if e != nil {
			fmt.Println(e)
			// t.Errorf("Error Found : ", e)
		}

		xGrabService.Log = logconf
		//===================================================================

		//===================================================================
		//==Data Setting and Destination Save =====================

		xGrabService.ServGrabber.DataSettings = make(map[string]*sdt.DataSetting)
		xGrabService.DestDbox = make(map[string]*sdt.DestInfo)

		// ==For Every Data Setting ===============================
		tempDataSetting := sdt.DataSetting{}
		tempDestInfo := sdt.DestInfo{}

		tempDataSetting.RowSelector = "#tab_conbox li:nth-child(1) .sjtable .listshuju tbody tr"
		tempDataSetting.Column(0, &sdt.GrabColumn{Alias: "Code", Selector: "td:nth-child(1)"})
		tempDataSetting.Column(0, &sdt.GrabColumn{Alias: "LongSpeculation", Selector: "td:nth-child(2)"})
		tempDataSetting.Column(0, &sdt.GrabColumn{Alias: "ShortSpeculation", Selector: "td:nth-child(3)"})

		xGrabService.ServGrabber.DataSettings["DATA01"] = &tempDataSetting //DATA01 use name in datasettings

		ci := dbox.ConnectionInfo{}
		ci.Host = "E:\\WORKS\\data_test\\vale\\Data_Grab.csv"
		ci.Database = ""
		ci.UserName = ""
		ci.Password = ""
		ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")

		tempDestInfo.Collection = ""
		tempDestInfo.Desttype = "csv"

		tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
		if e != nil {
			fmt.Println(e)
			// t.Errorf("Error Found : ", e)
		}

		xGrabService.DestDbox["DATA01"] = &tempDestInfo
		fmt.Println(xGrabService)
		//===================================================================

		e = xGrabService.StartService()
		if e != nil {
			fmt.Println(e)
			// t.Errorf("Error Found : ", e)
		} else {

			for i := 0; i < 100; i++ {
				fmt.Printf(".")
				time.Sleep(3000 * time.Millisecond)
			}

			e = xGrabService.StopService()
			if e != nil {
				fmt.Println(e)
				// t.Errorf("Error Found : ", e)
			}
		}
		return xGrabService
	*/
}
