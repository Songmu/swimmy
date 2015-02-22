package swimmy

import (
	"fmt"
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

func TestAgentCollectValues(t *testing.T) {
	a := NewAgent("./test", 3)
	v := a.collectValues()

	fmt.Printf("%+v\n", v)
	if v[0].service != "test1" {
		t.Errorf("something wrong")
	}

	if len(v) <= 0 {
		t.Errorf("something wrong")
	}

}
