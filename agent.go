package swimmy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
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

type postValue struct {
	service string
	values  []metricValue
}

func (a *agent) Watch() <-chan []postValue {
	ch := make(chan []postValue)
	timer := make(chan time.Time)

	go func() {
		next := time.Now()
		for {
			timer <- <-time.After(next.Sub(time.Now()))
			next = next.Add(time.Duration(a.interval) * time.Minute)
		}
	}()

	go func() {
		for _ = range timer {
			go func() {
				ch <- a.collectValues()
			}()
		}
	}()

	return ch
}

func (a *agent) collectValues() []postValue {
	result := []postValue{}

	resultChan := make(chan postValue, a.procs)

	go func() {
		var wg sync.WaitGroup
		for _, co := range a.collectors() {
			wg.Add(1)
			go func(co *collector) {
				defer wg.Done()

				resultChan <- postValue{
					values:  co.collectValues(),
					service: co.ServiceName(),
				}
			}(co)
		}
		wg.Wait()
		close(resultChan)
	}()

	for p := range resultChan {
		result = append(result, p)
	}

	return result
}
