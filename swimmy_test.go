package main

import (
	"regexp"
	"testing"
)

func TestCollectors(t *testing.T) {
	s := NewSwimmy(Args{
		Dir:      "./test",
		Interval: 5,
	})
	collectors := s.collectors()

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
	s := NewSwimmy(Args{
		Dir:      "./test",
		Interval: 3,
	})
	v := s.collectValues()

	if v[0].service != "test1" {
		t.Errorf("something wrong")
	}

	if len(v) <= 0 {
		t.Errorf("something wrong")
	}

}
