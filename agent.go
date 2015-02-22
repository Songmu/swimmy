package swimmy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type agent struct {
	dir      string
	procs    int
	interval uint
}

func NewAgent(dir string, interval uint) *agent {
	if interval == 0 {
		interval = 1
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Printf("Failed to create agent of dir:[%s] \"%s\"", dir, err)
		os.Exit(1)
	}

	_, err = ioutil.ReadDir(absDir)
	if err != nil {
		log.Printf("Can't read Directory : [%s]. %s\n", dir, err)
		os.Exit(1)
	}

	return &agent{
		dir:      absDir,
		interval: interval,
		procs:    1,
	}
}

func (a *agent) collectors() []*collector {
	collectors := []*collector{}

	fileInformations, err := ioutil.ReadDir(a.dir)
	if err != nil {
		fmt.Printf("Can't read Directory : [%s]. %s\n", a.dir, err)
		return collectors
	}

	for _, info := range fileInformations {
		name := info.Name()
		if !info.IsDir() || name[0] == '.' {
			continue
		}

		c := newCollector(filepath.Join(a.dir, name))
		collectors = append(collectors, c)
	}
	return collectors
}
