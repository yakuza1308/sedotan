package modules

import (
	"bufio"
	"fmt"
	"github.com/eaciit/cast"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	"github.com/eaciit/toolkit"
	"os"
	"strconv"
	"strings"
	"time"
)

type HistoryModule struct {
	filepathName, nameid, logPath string
	humanDate                     string
	rowgrabbed, rowsaved          float64
}

var (
	filepath = wd + "data\\History\\"
)

func NewHistory(nameid string) *HistoryModule {
	h := new(HistoryModule)

	dateNow := cast.Date2String(time.Now(), "YYYYMM") //time.Now()
	path := filepath + nameid + "-" + dateNow + ".csv"
	h.filepathName = path
	h.nameid = nameid
	return h
}

func (h *HistoryModule) OpenHistory() interface{} {
	var config = map[string]interface{}{"useheader": true, "delimiter": ",", "dateformat": "MM-dd-YYYY"}
	ci := &dbox.ConnectionInfo{h.filepathName, "", "", "", config}
	c, e := dbox.NewConnection("csv", ci)
	if e != nil {
		return e.Error()
	}

	e = c.Connect()
	if e != nil {
		return e.Error()
	}
	defer c.Close()

	csr, e := c.NewQuery().Select("*").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	if csr == nil {
		return "Cursor not initialized"
	}
	defer csr.Close()

	ds, e := csr.Fetch(nil, 0, false)
	if e != nil {
		return e.Error()
	}

	var history = []interface{}{} //toolkit.M{}
	for i, v := range ds.Data {
		// layout := "2006/01/02 15:04:05"
		castDate, _ := time.Parse(time.RFC3339, v.(toolkit.M)["grabdate"].(string))
		h.humanDate = cast.Date2String(castDate, "YYYY/MM/dd HH:mm:ss")
		h.rowgrabbed, _ = strconv.ParseFloat(v.(toolkit.M)["rowgrabbed"].(string), 64)
		h.rowsaved, _ = strconv.ParseFloat(v.(toolkit.M)["rowsaved"].(string), 64)

		var addToMap = toolkit.M{}
		addToMap.Set("id", i+1)
		addToMap.Set("datasettingname", v.(toolkit.M)["datasettingname"])
		addToMap.Set("grabdate", h.humanDate)
		addToMap.Set("grabstatus", v.(toolkit.M)["grabstatus"])
		addToMap.Set("rowgrabbed", h.rowgrabbed)
		addToMap.Set("rowsaved", h.rowsaved)
		addToMap.Set("notehistory", v.(toolkit.M)["note"])
		addToMap.Set("recfile", v.(toolkit.M).Get("recfile"))
		addToMap.Set("nameid", h.nameid)

		history = append(history, addToMap)
	}

	return history
}

func (h *HistoryModule) GetLog(datas []interface{}, date string) interface{} {
	for _, v := range datas {
		vMap, _ := toolkit.ToM(v)

		dateNow := time.Now()
		dateNowFormat := dateNow.Format(vMap["logconf"].(map[string]interface{})["filepattern"].(string))
		h.logPath = fmt.Sprintf("%s\\%s-%s", vMap["logconf"].(map[string]interface{})["logpath"], vMap["logconf"].(map[string]interface{})["filename"], dateNowFormat)
	}

	file, err := os.Open(h.logPath)
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	getHours := strings.Split(date, ":")
	containString := getHours[0] + ":" + getHours[1]
	scanner := bufio.NewScanner(file)
	lines := 0
	containLines := 0

	var logs []interface{}
	for scanner.Scan() {
		lines++
		contains := strings.Contains(scanner.Text(), containString)
		if contains {
			containLines = lines
		}

		if lines == containLines {
			logs = append(logs, "<li>"+scanner.Text()+"</li>")
		}
	}

	if err := scanner.Err(); err != nil {
		return err.Error()
	}

	var addSlice = toolkit.M{}
	addSlice.Set("logs", logs)

	return addSlice
}
