package sedotan

import (
	// "encoding/json"
	"fmt"
	"github.com/eaciit/toolkit"
	"testing"
)

func TestGrab(t *testing.T) {
	t.Skip()

	url := "http://www.ariefdarmawan.com"
	g := NewGrabber(url, "GET", &Config{})
	if e := g.Grab(nil); e != nil {
		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
		return
	}

	fmt.Printf("Result:\n%s\n", g.ResultString()[:200])
}

// func TestQuery(t *testing.T) {
// 	url := "http://www.ariefdarmawan.com"

// 	g := NewGrabber(url, "GET", nil)
// 	g.Config.DataSettings = make(map[string]*DataSetting)

// 	tempDataSetting := DataSetting{}
// 	tempDataSetting.RowSelector = "article"
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "Title", Selector: "h1.entry-title"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "Excerpt", Selector: ".entry-content"})

// 	g.Config.DataSettings["SELECT01"] = &tempDataSetting

// 	if e := g.Grab(nil); e != nil {
// 		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
// 		return
// 	}

// 	docs := []struct {
// 		Title   string
// 		Excerpt string
// 	}{}

// 	e := g.ResultFromHtml("SELECT01", &docs)
// 	if e != nil {
// 		t.Errorf("Unable to read: %s", e.Error())
// 	}
// 	fmt.Printf("Result:\n%s\n", func() string {
// 		ret := ""
// 		for _, doc := range docs {
// 			ret += "# " + doc.Title + "\n" +
// 				doc.Excerpt + "\n" +
// 				"================================================================" +
// 				"\n"
// 		}
// 		return ret
// 	}())
// }

// func TestQuery(t *testing.T) {
// 	url := "http://www.shfe.com.cn/en/products/Gold/"

// 	g := NewGrabber(url, "GET", nil)
// 	g.Config.DataSettings = make(map[string]*DataSetting)

// 	tempDataSetting := DataSetting{}
// 	tempDataSetting.RowSelector = "#tab_conbox li:nth-child(2) .sjtable .listshuju tbody tr"
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "Code", Selector: "td:nth-child(1)"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "LongSpeculation", Selector: "td:nth-child(2)"})
// 	tempDataSetting.Column(0, &GrabColumn{Alias: "ShortSpeculation", Selector: "td:nth-child(3)"})

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

func TestPost(t *testing.T) {
	url := "http://www.dce.com.cn/PublicWeb/MainServlet"
	GrabConfig := Config{}

	dataurl := toolkit.M{}
	dataurl["Pu00231_Input.trade_date"] = "20151214"
	dataurl["Pu00231_Input.variety"] = "i"
	dataurl["Pu00231_Input.trade_type"] = "0"
	dataurl["Submit"] = "Go"
	dataurl["action"] = "Pu00231_result"

	// postdata := toolkit.M{}.Set("formvalues", dataurl)
	// fmt.Println(postdata)

	// str := `{
	//            "formvalues":
	//              {
	//                "Pu00231_Input.trade_date": "20151214",
	//                "Pu00231_Input.variety": "i",
	//                "Pu00231_Input.trade_type": "0",
	//                "Submit": "Go",
	//                "action": "Pu00231_result"
	//              }
	//          }`
	// res := toolkit.M{}
	// json.Unmarshal([]byte(str), &res)

	GrabConfig.setFormValues(dataurl)
	g := NewGrabber(url, "POST", &GrabConfig)
	// fmt.Println(g.Config.PostData)
	g.Config.DataSettings = make(map[string]*DataSetting)

	tempDataSetting := DataSetting{}
	tempDataSetting.RowSelector = "table .table tbody tr"
	tempDataSetting.Column(0, &GrabColumn{Alias: "Contract", Selector: "td:nth-child(1)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "Open", Selector: "td:nth-child(2)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "High", Selector: "td:nth-child(3)"})

	g.Config.DataSettings["SELECT01"] = &tempDataSetting

	if e := g.Grab(nil); e != nil {
		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
		return
	}

	docs := []toolkit.M{}

	e := g.ResultFromHtml("SELECT01", &docs)
	if e != nil {
		t.Errorf("Unable to read: %s", e.Error())
	}

	for _, doc := range docs {
		fmt.Println(doc)
	}

	// dataurl := url.Values{}
	// dataurl.Add("Pu00231_Input.trade_date", "20151214")
	// dataurl.Add("Pu00231_Input.variety", "i")
	// dataurl.Add("Pu00231_Input.trade_type", 0)
	// dataurl.Add("Submit", "Go")
	// dataurl.Add("action", "Pu00231_result")

	// fmt.Println(dataurl)
	// fmt.Println("LINE 148", g.Config.PostData)
	// fmt.Println(g.ResultString())
}
