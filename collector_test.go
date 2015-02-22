package swimmy

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
)

func TestCollectFromCmd(t *testing.T) {

	c, _ := newCollector("./test/test1")

	if c.ServiceName() != "test1" {
		t.Errorf("ServiceName should be equals directory name")
	}
	cmd, _ := filepath.Abs("./test/test1/foo")
	r, _ := c.collectFromCmd(cmd)

	if r["foo"] != 1 {
		t.Errorf("foo is not collected")
	}

	if r["foo.bar"] != 15.5 {
		t.Errorf("foo.bar is not collected")
	}
}

func TestCollectValues(t *testing.T) {
	c, _ := newCollector("./test/test1")
	v, _ := c.collectValues()

	r := make(map[string]float64)
	for _, mv := range v {
		r[mv.Name] = mv.Value
	}

	if r["foo"] != 1 {
		t.Errorf("foo is not collected")
	}

	if r["foo.bar"] != 15.5 {
		t.Errorf("foo.bar is not collected")
	}

	if r["bar"] != 18.8 {
		t.Errorf("bar is not collected")
	}

	if r["baz.hoge"] != 33 {
		t.Errorf("baz.hoge is not collected")
	}

	if r["baz.hoge.dummy"] != 15 {
		t.Errorf("baz.hoge is not collected")
	}

	valuesJSON, _ := json.Marshal(v)
	fmt.Printf("%s\n", valuesJSON)
	// [{"time":1424588215,"name":"bar","value":18.8},{"time":1424588215,"name":"baz.hoge","value":33},{"time":1424588215,"name":"baz.hoge.dummy","value":15},{"time":1424588215,"name":"foo","value":1},{"time":1424588215,"name":"foo.bar","value":15.5}]
}
