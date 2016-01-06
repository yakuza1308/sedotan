package modules

import (
	"bufio"
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	"github.com/eaciit/toolkit"
	"os"
	"strings"
	"time"
)

type HistoryModule struct {
	filepathName, nameid, logPath string
}

var (
	filepath = wd + "data\\History\\"
)

func NewHistory(nameid string) *HistoryModule {
	h := new(HistoryModule)

	dateNow := time.Now()
	path := filepath + nameid + "-" + dateNow.Format("200601") + ".csv"
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
	for _, v := range ds.Data {
		var addToMap = toolkit.M{}
		addToMap.Set("id", v.(toolkit.M)["id"])
		addToMap.Set("datasettingname", v.(toolkit.M)["datasettingname"])
		addToMap.Set("grabdate", v.(toolkit.M)["grabdate"])
		addToMap.Set("grabstatus", v.(toolkit.M)["grabstatus"])
		addToMap.Set("rowgrabbed", v.(toolkit.M)["rowgrabbed"])
		addToMap.Set("rowsaved", v.(toolkit.M)["rowsaved"])
		addToMap.Set("notehistory", v.(toolkit.M)["notehistory"])
		addToMap.Set("nameid", h.nameid)

		history = append(history, addToMap)
	}

	fmt.Sprintf("history=%v\n", history)
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
