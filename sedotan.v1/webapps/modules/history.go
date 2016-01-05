package modules

import (
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	// "io"
	// "os"
	"time"
)

type HistoryModule struct {
	filepathName string
}

var (
	filepath = wd + "data\\History\\"
)

func NewHistory(nameid string) *HistoryModule {
	h := new(HistoryModule)

	dateNow := time.Now()
	path := filepath + nameid + "-" + dateNow.Format("200601") + ".csv"
	h.filepathName = path
	return h
}

func (h *HistoryModule) OpenHistory() interface{} {
	// buka file

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

	// file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	// if err != nil {
	// 	fmt.Printf("err1:%v\n", err)
	// 	return err
	// }
	// // checkError(err)
	// defer file.Close()

	// // baca file
	// var text = make([]byte, 1024)
	// fmt.Printf("text1:%v\n", text)
	// for {
	// 	n, err := file.Read(text)
	// 	if err != io.EOF {
	// 		fmt.Printf("err2:%v\n", err)
	// 		return err
	// 	}
	// 	if n == 0 {
	// 		break
	// 	}
	// }
	fmt.Sprintf("text:%v\n", ds.Data)
	// checkError(err)
	return ds.Data
}
