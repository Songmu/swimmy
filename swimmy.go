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

type Swimmy struct {
	dir      string
	procs    int
	interval uint
}

func NewSwimmy(dir string, interval uint) *Swimmy {
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

	return &Swimmy{
		dir:      absDir,
		interval: interval,
		procs:    1,
	}
}

func (s *Swimmy) collectors() []*collector {
	collectors := []*collector{}

	fileInformations, err := ioutil.ReadDir(s.dir)
	if err != nil {
		fmt.Printf("Can't read Directory : [%s]. %s\n", s.dir, err)
		return collectors
	}

	for _, info := range fileInformations {
		name := info.Name()
		if !info.IsDir() || name[0] == '.' {
			continue
		}

		c := newCollector(filepath.Join(s.dir, name))
		collectors = append(collectors, c)
	}
	return collectors
}

type postValue struct {
	service  string
	values   []metricValue
	retryCnt uint
}

func (s *Swimmy) watch() chan *postValue {
	ch := make(chan *postValue)
	timer := make(chan time.Time)

	go func() {
		next := time.Now()
		for {
			timer <- <-time.After(next.Sub(time.Now()))
			next = next.Add(time.Duration(s.interval) * time.Minute)
		}
	}()

	go func() {
		for _ = range timer {
			go func() {
				values := s.collectValues()
				for _, v := range values {
					ch <- v
				}
			}()
		}
	}()

	return ch
}

func (s *Swimmy) collectValues() []*postValue {
	result := []*postValue{}

	resultChan := make(chan *postValue, s.procs)

	go func() {
		var wg sync.WaitGroup
		for _, co := range s.collectors() {
			wg.Add(1)
			go func(co *collector) {
				defer wg.Done()

				resultChan <- &postValue{
					values:   co.collectValues(),
					service:  co.ServiceName(),
					retryCnt: 0,
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
