package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type metricValue struct {
	Time  int64   `json:"time"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type collector struct {
	dir   string
	procs int
}

func newCollector(dir string) *collector {
	// dir must be absolute path
	return &collector{
		dir:   dir,
		procs: 1,
	}
}

func (c *collector) ServiceName() string {
	return filepath.Base(c.dir)
}

func (c *collector) collectValues() []metricValue {
	var result []metricValue

	ch := c.gatherExecutable()
	resultChan := make(chan map[string]float64, c.procs)
	go func() {
		var wg sync.WaitGroup
		for path := range ch {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				v, err := c.collectFromCmd(path)
				if err != nil {
					log.Warn("error")
					return
				}
				resultChan <- v
			}(path)
		}
		wg.Wait()
		close(resultChan)
	}()

	time := time.Now().Unix()
	for r := range resultChan {
		for k, v := range r {
			result = append(result, metricValue{
				Time:  time,
				Name:  k,
				Value: v,
			})
		}
	}
	return result
}

func (c *collector) gatherExecutable() <-chan string {
	ch := make(chan string, c.procs)

	go func() {
		filepath.Walk(c.dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if (info.Name())[0] == '.' {
					return filepath.SkipDir
				}
				return nil
			}

			if (info.Mode() & 0111) != 0 {
				ch <- path
			}
			return nil
		})
		close(ch)
	}()

	return ch
}

func (c *collector) collectFromCmd(cmd string) (map[string]float64, error) {
	log.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Debug("Executing command")

	// os.Setenv(pluginConfigurationEnvName, "")
	stdout, stderr, err := runCommand(cmd)

	if err != nil {
		log.WithFields(logrus.Fields{
			"cmd":    cmd,
			"err":    err,
			"stdout": stdout,
			"stderr": stderr,
		}).Warn("Failed to execute command (skip these metrics)")
		return nil, err
	}

	rel, err := filepath.Rel(c.dir, cmd)
	if err != nil {
		log.WithFields(logrus.Fields{
			"cmd": cmd,
			"err": err,
		}).Warn("Failed to resolve relative path")
		return nil, err
	}
	baseKey := strings.Replace(rel, string(filepath.Separator), ".", -1)

	results := make(map[string]float64)
	for _, line := range strings.Split(stdout, "\n") {
		if line == "" {
			continue
		}
		// Key, value or value only
		items := strings.Split(line, "\t")
		l := len(items)
		if l > 2 {
			continue
		}

		k := baseKey
		vIdx := 0
		if l == 2 {
			vIdx = 1
			k = baseKey + "." + items[0]
		}

		v, err := strconv.ParseFloat(items[vIdx], 64)
		if err != nil {
			log.WithFields(logrus.Fields{
				"err":   err,
				"value": items[vIdx],
			}).Warn("Failed to parse value.")
			continue
		}
		results[k] = v
	}

	return results, nil
}

// runCommand runs command (in one string) and returns stdout, stderr strings.
func runCommand(command string) (string, string, error) {
	var outBuffer, errBuffer bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()

	if err != nil {
		return "", "", err
	}

	return string(outBuffer.Bytes()), string(errBuffer.Bytes()), nil
}
