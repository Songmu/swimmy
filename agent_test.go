package swimmy

import (
	"regexp"
	"testing"
)

func TestCollectors(t *testing.T) {
	a := NewAgent("./test", 5)

	collectors := a.collectors()

	var co *collector
	for _, c := range collectors {
		if c.ServiceName() == "test1" {
			co = c
			break
		}
	}

	re, _ := regexp.Compile("/test/test1$")
	if !re.MatchString(co.dir) {
		t.Errorf("something wrong")
	}

}
