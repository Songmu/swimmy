package swimmy

import (
	"fmt"
	"os"
	"testing"
)

func TestCollectFromCmd(t *testing.T) {

	c, _ := newCollector(".")
	r, _ := c.collectFromCmd(os.Getenv("GOPATH") + "/src/github.com/Songmu/swimmy/test/foo.go")
	fmt.Printf("%+v\n", r)

	if r["test.foo.go"] != 1 {
		t.Errorf("test.foo.go is not collected")
	}

	if r["test.foo.go.sample"] != 15.5 {
		t.Errorf("test.foo.go.cample is not collected")
	}
}

func TestCollectValues(t *testing.T) {
	c, _ := newCollector(".")

	c.collectValues()
}
