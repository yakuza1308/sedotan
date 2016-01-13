package sedotan

import (
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"testing"
	"time"
)

// func TestGrab(t *testing.T) {
// 	t.Skip()

// 	url := "http://www.ariefdarmawan.com"
// 	g := NewGrabber(url, "GET", &Config{})
// 	if e := g.Grab(nil); e != nil {
// 		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
// 		return
// 	}

// 	fmt.Printf("Result:\n%s\n", g.ResultString()[:200])
// }

// func TestPost(t *testing.T) {
// 	url := "http://www.dce.com.cn/PublicWeb/MainServlet"
// 	GrabConfig := Config{}

// 	dataurl := toolkit.M{}
// 	dataurl["Pu00231_Input.trade_date"] = "20151214"
// 	dataurl["Pu00231_Input.variety"] = "i"
// 	dataurl["Pu00231_Input.trade_type"] = "0"
// 	dataurl["Submit"] = "Go"
// 	dataurl["action"] = "Pu00231_result"

// 	GrabConfig.SetFormValues(dataurl)
// 	g := NewGrabber(url, "POST", &GrabConfig)

// 	g.Config.DataSettings = make(map[string]*DataSetting)

// 	tempDataSetting := DataSetting{}
// 	tempDataSetting.RowSelector = "table .table tbody tr"
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "Contract", Selector: "td:nth-child(1)"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "Open", Selector: "td:nth-child(2)"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "High", Selector: "td:nth-child(3)"})

// 	g.Config.DataSettings["SELECT01"] = &tempDataSetting

// 	if e := g.Grab(nil); e != nil {
// 		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
// 		return
// 	}

// 	docs := []toolkit.M{}

// 	e := g.ResultFromHtml("SELECT01", &docs)
// 	if e != nil {
// 		t.Errorf("Unable to read: %s", e.Error())
// 	}

// 	for _, doc := range docs {
// 		fmt.Println(doc)
// 	}

// }

// func TestServiceGrabGet(t *testing.T) {

// 	xGrabService := NewGrabService()
// 	xGrabService.Name = "getgoldshfecom"
// 	xGrabService.Url = "http://www.shfe.com.cn/en/products/Gold/"

// 	xGrabService.SourceType = SourceType_Http

// 	xGrabService.GrabInterval = 2 * time.Minute
// 	xGrabService.TimeOutInterval = 5 * time.Second //time.Hour, time.Minute, time.Second

// 	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, "seconds")

// 	//==For Data Grab Config/Data Grabber          ===========================================
// 	// tempGrab.ServGrabber = sedotan.NewGrabber(tempGrab.Url, mapVal.Get("calltype", "").(string), &GrabConfig)

// 	grabConfig := Config{}
// 	// if has grabconfig

// 	// CallType     string
// 	// FormValues   toolkit.M
// 	// AuthType     string
// 	// AuthUserId   string
// 	// AuthPassword string
// 	//==================

// 	xGrabService.ServGrabber = NewGrabber(xGrabService.Url, "GET", &grabConfig)

// 	//===================================================================

// 	//==For Data Log          ===========================================

// 	// 	logpath = tempLogConf.Get("logpath", "").(string)
// 	// 	filename = tempLogConf.Get("filename", "").(string)
// 	// 	filepattern = tempLogConf.Get("filepattern", "").(string)

// 	logpath := "E:\\data\\vale\\log"
// 	filename := "LOG-GRABSHFETEST-%s"
// 	filepattern := "YYYYMMDD"

// 	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
// 	if e != nil {
// 		t.Errorf("Error Found : ", e)
// 	}

// 	xGrabService.Log = logconf
// 	//===================================================================

// 	//===================================================================
// 	//==Data Setting and Destination Save =====================

// 	xGrabService.ServGrabber.DataSettings = make(map[string]*DataSetting)
// 	xGrabService.DestDbox = make(map[string]*DestInfo)

// 	// ==For Every Data Setting ===============================
// 	tempDataSetting := DataSetting{}
// 	tempDestInfo := DestInfo{}

// 	tempDataSetting.RowSelector = "#tab_conbox li:nth-child(2) .sjtable .listshuju tbody tr"
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "Code", Selector: "td:nth-child(1)"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "LongSpeculation", Selector: "td:nth-child(2)"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "ShortSpeculation", Selector: "td:nth-child(3)"})

// 	xGrabService.ServGrabber.DataSettings["DATA01"] = &tempDataSetting //DATA01 use name in datasettings

// 	ci := dbox.ConnectionInfo{}
// 	ci.Host = "E:\\data\\vale\\Data_Grab.csv"
// 	ci.Database = ""
// 	ci.UserName = ""
// 	ci.Password = ""
// 	ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")

// 	tempDestInfo.Collection = ""
// 	tempDestInfo.Desttype = "csv"

// 	tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
// 	if e != nil {
// 		t.Errorf("Error Found : ", e)
// 	}

// 	xGrabService.DestDbox["DATA01"] = &tempDestInfo
// 	//=History===========================================================
// 	xGrabService.HistoryPath = "E:\\data\\vale\\history\\"
// 	xGrabService.HistoryRecPath = "E:\\data\\vale\\historyrec\\"

// 	// xGrabService.HistDbox = &tempHistInfo
// 	//===================================================================

// 	e = xGrabService.StartService()
// 	if e != nil {
// 		t.Errorf("Error Found : ", e)
// 	} else {

// 		fmt.Printf("[SUM] start %s, grab %d times, data retreive %d rows, error %d times\n", xGrabService.StartDate, xGrabService.GrabCount, xGrabService.RowGrabbed, xGrabService.ErrorFound)
// 		for i := 0; i < 100; i++ {
// 			fmt.Printf(".")
// 			time.Sleep(1000 * time.Millisecond)
// 		}

// 		fmt.Println()
// 		fmt.Printf("[SUM] start %s, grab %d times, data retreive %d rows, error %d times\n", xGrabService.StartDate, xGrabService.GrabCount, xGrabService.RowGrabbed, xGrabService.ErrorFound)

// 		for i := 0; i < 100; i++ {
// 			fmt.Printf(".")
// 			time.Sleep(1000 * time.Millisecond)
// 		}

// 		fmt.Println()
// 		fmt.Printf("[SUM] start %s, grab %d times, data retreive %d rows, error %d times\n", xGrabService.StartDate, xGrabService.GrabCount, xGrabService.RowGrabbed, xGrabService.ErrorFound)

// 		e = xGrabService.StopService()
// 		if e != nil {
// 			t.Errorf("Error Found : ", e)
// 		}
// 	}
// }

func TestServiceGrabPost(t *testing.T) {

	xGrabService := NewGrabService()
	xGrabService.Name = "irondcecom"
	xGrabService.Url = "http://www.dce.com.cn/PublicWeb/MainServlet"

	xGrabService.SourceType = SourceType_HttpHtml

	xGrabService.GrabInterval = 5 * time.Minute
	xGrabService.TimeOutInterval = 10 * time.Second //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, "seconds")

	//==For Data Grab Config/Data Grabber          ===========================================
	// tempGrab.ServGrabber = sedotan.NewGrabber(tempGrab.Url, mapVal.Get("calltype", "").(string), &GrabConfig)

	grabConfig := Config{}

	dataurl := toolkit.M{}
	dataurl["Pu00231_Input.trade_date"] = "20151214"
	dataurl["Pu00231_Input.variety"] = "i"
	dataurl["Pu00231_Input.trade_type"] = "0"
	dataurl["Submit"] = "Go"
	dataurl["action"] = "Pu00231_result"

	grabConfig.SetFormValues(dataurl)

	// if has grabconfig

	// CallType     string
	// FormValues   toolkit.M
	// AuthType     string
	// AuthUserId   string
	// AuthPassword string
	//==================

	xGrabService.ServGrabber = NewGrabber(xGrabService.Url, "POST", &grabConfig)

	//===================================================================

	//==For Data Log          ===========================================

	// 	logpath = tempLogConf.Get("logpath", "").(string)
	// 	filename = tempLogConf.Get("filename", "").(string)
	// 	filepattern = tempLogConf.Get("filepattern", "").(string)

	logpath := "E:\\data\\vale\\log"
	filename := "LOG-GRABDCETEST-%s"
	filepattern := "20060102"

	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.Log = logconf
	//===================================================================

	//===================================================================
	//==Data Setting and Destination Save =====================

	xGrabService.ServGrabber.DataSettings = make(map[string]*DataSetting)
	xGrabService.DestDbox = make(map[string]*DestInfo)

	// ==For Every Data Setting ===============================
	tempDataSetting := DataSetting{}
	tempDestInfo := DestInfo{}

	tempDataSetting.RowSelector = "table .table tbody tr"
	tempDataSetting.Column(0, &GrabColumn{Alias: "Contract", Selector: "td:nth-child(1)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "Open", Selector: "td:nth-child(2)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "High", Selector: "td:nth-child(3)"})

	andCondition := []interface{}{}
	andCondition = append(andCondition, map[string]interface{}{"Contract": map[string]interface{}{"$ne": "Contract"}})
	andCondition = append(andCondition, map[string]interface{}{"Contract": map[string]interface{}{"$ne": "Iron Ore Subtotal"}})
	andCondition = append(andCondition, map[string]interface{}{"Contract": map[string]interface{}{"$ne": "Total"}})

	// orCondition[0] = map[string]string{"Contract": "Contract"}
	// orCondition[1] = map[string]string{"Contract": "Iron Ore Subtotal"}
	// orCondition[2] = map[string]string{"Contract": "Total"}

	tempFilterCond := toolkit.M{}.Set("$and", andCondition)
	tempDataSetting.SetFilterCond(tempFilterCond)
	// -Check "filtercond" in config-
	// tempFilterCond, e = toolkit.ToM(mapxVal.Get("filtercond", nil).(map[string]interface{}))
	// tempDataSetting.SetFilterCond(tempFilterCond)

	xGrabService.ServGrabber.DataSettings["DATA01"] = &tempDataSetting //DATA01 use name in datasettings

	ci := dbox.ConnectionInfo{}
	ci.Host = "E:\\data\\vale\\Data_GrabIronTest.csv"
	ci.Database = ""
	ci.UserName = ""
	ci.Password = ""
	ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")

	tempDestInfo.Collection = ""
	tempDestInfo.Desttype = "csv"

	tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.DestDbox["DATA01"] = &tempDestInfo
	//=History===========================================================
	xGrabService.HistoryPath = "E:\\data\\vale\\history\\"
	xGrabService.HistoryRecPath = "E:\\data\\vale\\historyrec\\"
	//===================================================================

	e = xGrabService.StartService()
	if e != nil {
		t.Errorf("Error Found : ", e)
	} else {

		for i := 0; i < 100; i++ {
			fmt.Printf(".")
			time.Sleep(1000 * time.Millisecond)
		}

		e = xGrabService.StopService()
		if e != nil {
			t.Errorf("Error Found : ", e)
		}
	}
}

func TestServiceGrabLoggin(t *testing.T) {

	xGrabService := NewGrabService()
	xGrabService.Name = "localtest"
	xGrabService.Url = "http://localhost:8000"

	xGrabService.SourceType = SourceType_HttpHtml

	xGrabService.GrabInterval = 5 * time.Minute
	xGrabService.TimeOutInterval = 10 * time.Second //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, "seconds")

	//==For Data Grab Config/Data Grabber          ===========================================
	// tempGrab.ServGrabber = sedotan.NewGrabber(tempGrab.Url, mapVal.Get("calltype", "").(string), &GrabConfig)

	grabConfig := Config{}

	grabConfig.AuthType = "session"
	// grabConfig.AuthUserId = mapValConfig.Get("authuserid", "").(string)
	// grabConfig.AuthPassword = mapValConfig.Get("authpassword", "").(string)

	grabConfig.LoginUrl = "http://localhost:8000/login"
	grabConfig.LogoutUrl = "http://localhost:8000/logout"
	grabConfig.LoginValues = toolkit.M{}.Set("name", "alip").Set("password", "test")

	// if has grabconfig

	// CallType     string
	// FormValues   toolkit.M
	// AuthType     string
	// AuthUserId   string
	// AuthPassword string
	//==================

	xGrabService.ServGrabber = NewGrabber(xGrabService.Url, "POST", &grabConfig)

	//===================================================================

	//==For Data Log          ===========================================

	// 	logpath = tempLogConf.Get("logpath", "").(string)
	// 	filename = tempLogConf.Get("filename", "").(string)
	// 	filepattern = tempLogConf.Get("filepattern", "").(string)

	logpath := "E:\\data\\vale\\log"
	filename := "LOG-LOCALTEST-%s"
	filepattern := "20060102"

	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.Log = logconf
	//===================================================================

	//===================================================================
	//==Data Setting and Destination Save =====================

	xGrabService.ServGrabber.DataSettings = make(map[string]*DataSetting)
	xGrabService.DestDbox = make(map[string]*DestInfo)

	// ==For Every Data Setting ===============================
	tempDataSetting := DataSetting{}
	tempDestInfo := DestInfo{}

	tempDataSetting.RowSelector = "table tr"
	tempDataSetting.Column(0, &GrabColumn{Alias: "Number", Selector: "td:nth-child(1)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "FirstName", Selector: "td:nth-child(2)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "LastName", Selector: "td:nth-child(3)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "Points", Selector: "td:nth-child(4)"})

	tempFilterCond := toolkit.M{}.Set("FirstName", map[string]interface{}{"$ne": ""})
	tempDataSetting.SetFilterCond(tempFilterCond)
	// -Check "filtercond" in config-
	// tempFilterCond, e = toolkit.ToM(mapxVal.Get("filtercond", nil).(map[string]interface{}))
	// tempDataSetting.SetFilterCond(tempFilterCond)

	xGrabService.ServGrabber.DataSettings["DATA01"] = &tempDataSetting //DATA01 use name in datasettings

	ci := dbox.ConnectionInfo{}
	ci.Host = "E:\\data\\vale\\Data_GrabLocal.csv"
	ci.Database = ""
	ci.UserName = ""
	ci.Password = ""
	ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")

	tempDestInfo.Collection = ""
	tempDestInfo.Desttype = "csv"

	tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.DestDbox["DATA01"] = &tempDestInfo
	//=History===========================================================
	xGrabService.HistoryPath = "E:\\data\\vale\\history\\"
	xGrabService.HistoryRecPath = "E:\\data\\vale\\historyrec\\"
	//===================================================================

	e = xGrabService.StartService()
	if e != nil {
		t.Errorf("Error Found : ", e)
	} else {

		for i := 0; i < 100; i++ {
			fmt.Printf(".")
			time.Sleep(1000 * time.Millisecond)
		}

		e = xGrabService.StopService()
		if e != nil {
			t.Errorf("Error Found : ", e)
		}
	}
}

func TestServiceGrabDocument(t *testing.T) {
	var e error

	xGrabService := NewGrabService()
	xGrabService.Name = "iopriceindices"

	xGrabService.SourceType = SourceType_DocExcel

	xGrabService.GrabInterval = 5 * time.Minute
	xGrabService.TimeOutInterval = 10 * time.Second //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, "seconds")

	//==must have grabconf and Connection info inside grabconf         ===========================================
	// mapValConfig, e := toolkit.ToM(mapVal.Get("grabconf", nil).(map[string]interface{}))
	// mapConnVal, e := toolkit.ToM(mapValConfig.Get("connectioninfo", nil).(map[string]interface{}))

	ci := dbox.ConnectionInfo{}

	ci.Host = "E:\\data\\sample\\IO Price Indices.xlsm"
	// ci.Database = mapConnVal.Get("database", "").(string)
	// ci.UserName = mapConnVal.Get("userName", "").(string)
	// ci.Password = mapConnVal.Get("password", "").(string)

	// if have setting inside of connection info
	// ci.Settings, e = toolkit.ToM(tempSetting.(map[string]interface{}))
	ci.Settings = nil

	xGrabService.ServGetData, e = NewGetDatabase(ci.Host, "xlsx", &ci)

	//===================================================================

	//==For Data Log          ===========================================

	// 	logpath = tempLogConf.Get("logpath", "").(string)
	// 	filename = tempLogConf.Get("filename", "").(string)
	// 	filepattern = tempLogConf.Get("filepattern", "").(string)

	logpath := "E:\\data\\vale\\log"
	filename := "LOG-LOCALXLSX-%s"
	filepattern := "20060102"

	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.Log = logconf
	//===================================================================

	//===================================================================
	//==Data Setting and Destination Save =====================

	xGrabService.ServGetData.CollectionSettings = make(map[string]*CollectionSetting)
	xGrabService.DestDbox = make(map[string]*DestInfo)

	// ==For Every Data Setting ===============================
	tempDataSetting := CollectionSetting{}
	tempDestInfo := DestInfo{}

	// .Collection = mapxVal.Get("rowselector", "").(string)
	tempDataSetting.Collection = "HIST"
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "Date", Selector: "1"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "Platts 62% Fe IODEX", Selector: "2"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "Platts 65% Fe", Selector: "4"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "TSI 62% Fe", Selector: "15"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "TSI 65% Fe", Selector: "16"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "TSI 62% Fe LOW ALUMINA", Selector: "17"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "MB 62% Fe", Selector: "26"})
	tempDataSetting.SelectColumn = append(tempDataSetting.SelectColumn, &GrabColumn{Alias: "MB 65% Fe", Selector: "29"})

	// tempDataSetting.SetFilterCond(tempFilterCond)
	// -Check "filtercond" in config-
	// tempFilterCond, e = toolkit.ToM(mapxVal.Get("filtercond", nil).(map[string]interface{}))
	// tempDataSetting.SetFilterCond(tempFilterCond)

	xGrabService.ServGetData.CollectionSettings["DATA01"] = &tempDataSetting //DATA01 use name in datasettings

	ci = dbox.ConnectionInfo{}
	ci.Host = "localhost:27017"
	ci.Database = "valegrab"
	ci.UserName = ""
	ci.Password = ""
	// ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")
	// setting will be depend on config file

	tempDestInfo.Collection = "iopriceindices"
	tempDestInfo.Desttype = "mongo"

	tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.Desttype, &ci)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.DestDbox["DATA01"] = &tempDestInfo
	//=History===========================================================
	xGrabService.HistoryPath = "E:\\data\\vale\\history\\"
	xGrabService.HistoryRecPath = "E:\\data\\vale\\historyrec\\"
	//===================================================================

	e = xGrabService.StartService()
	if e != nil {
		t.Errorf("Error Found : ", e)
	} else {

		for i := 0; i < 100; i++ {
			fmt.Printf(".")
			time.Sleep(1000 * time.Millisecond)
		}

		e = xGrabService.StopService()
		if e != nil {
			t.Errorf("Error Found : ", e)
		}
	}
}
