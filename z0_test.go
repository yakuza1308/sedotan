package sedotan

import (
	"fmt"
	"testing"
)

func TestGrab(t *testing.T) {
	url := "http://www.ariefdarmawan.com"
	g := NewGrabber(url, "GET", nil)
	if e := g.Grab(nil); e != nil {
		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
	}

	fmt.Printf("Result:\n%s\n", g.ResultString()[:200])
}
