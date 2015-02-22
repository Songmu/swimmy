package swimmy

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCollectFromCmd(t *testing.T) {

	c, _ := newCollector("./test")

	if c.ServiceName() != "test" {
		t.Errorf("ServiceName should be equals directory name")
	}

	r, _ := c.collectFromCmd(os.Getenv("GOPATH") + "/src/github.com/Songmu/swimmy/test/foo.go")
	fmt.Printf("%+v\n", r)

	if r["foo.go"] != 1 {
		t.Errorf("foo.go is not collected")
	}

	if r["foo.go.sample"] != 15.5 {
		t.Errorf("foo.go.cample is not collected")
	}
}

func TestCollectValues(t *testing.T) {
	c, _ := newCollector("./test")
	v, _ := c.collectValues()

	fmt.Printf("%+v\n", v)

	valuesJSON, _ := json.Marshal(v)
	fmt.Printf("%s\n", valuesJSON)
}
