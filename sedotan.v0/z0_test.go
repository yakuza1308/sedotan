package sedotan

import (
	"fmt"
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

func TestJQuery(t *testing.T) {
	url := "http://www.ariefdarmawan.com"

	g := NewGrabber(url, "GET", nil)
	g.RowSelector = "article"
	g.Column(0, &GrabColumn{Alias: "Title", Selector: "h1.entry-title"})
	g.Column(0, &GrabColumn{Alias: "Excerpt", Selector: ".entry-content"})
	if e := g.Grab(nil); e != nil {
		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
		return
	}

	docs := []struct {
		Title   string
		Excerpt string
	}{}

	e := g.ResultFromHtml(&docs)
	if e != nil {
		t.Errorf("Unable to read: %s", e.Error())
	}
	fmt.Printf("Result:\n%s\n", func() string {
		ret := ""
		for _, doc := range docs {
			ret += "# " + doc.Title + "\n" +
				doc.Excerpt + "\n" +
				"================================================================" +
				"\n"
		}
		return ret
	}())
}
