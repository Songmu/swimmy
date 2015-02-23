package main

import (
	"regexp"
	"testing"
)

func TestCollectors(t *testing.T) {
	s := newSwimmy(args{
		dir:      "./test",
		interval: 5,
	})
	collectors := s.collectors()

	co := collectors[0]

	re, _ := regexp.Compile("/test/test1$")
	if !re.MatchString(co.dir) {
		t.Errorf("something wrong")
	}

}

func TestAgentCollectValues(t *testing.T) {
	s := newSwimmy(args{
		service:  "hoge",
		dir:      "./test",
		interval: 3,
	})
	v := s.collectValues()

	if v[0].service != "hoge" {
		t.Errorf("something wrong")
	}

	if len(v) <= 0 {
		t.Errorf("something wrong")
	}

}
