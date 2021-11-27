package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	mffp "github.com/icemarkom/mffprober"
)

// PollFan executes HTTP POST to query the fan status, reporting error.
func PollFan(cfg *mffp.Config) (*mffp.FanData, error) {
	fd := new(mffp.FanData)

	url := fmt.Sprintf("http://%s/%s", cfg.Host, mffp.MFFTarget)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(mffp.MFFQuery)))
	if err != nil {
		return nil, fmt.Errorf("error formatting request: %w", err)
	}

	hc := &http.Client{
		Timeout: cfg.Timeout,
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error polling the fan: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, &fd)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}

	return fd, nil
}

// ProbeFan runs continuous fan poller.
func ProbeFan(cfg *mffp.Config) (int, error) {
	var (
		probeCount, failCount int
		err                   error
		fd                    *mffp.FanData
	)

	for probeCount = 1; cfg.Count == 0 || probeCount <= cfg.Count; probeCount++ {
		fd, err = PollFan(cfg)
		if err != nil {
			failCount++
			log.SetOutput(os.Stderr)
			err = fmt.Errorf("error reading fan information from %q: %w", cfg.Host, err)
			if !cfg.Quiet {
				log.Printf("Probe #%d (%d success, %d fail): %v.", probeCount, probeCount-failCount, failCount, err)
			}
			if cfg.MaxFailCount > 0 && failCount >= cfg.MaxFailCount {
				break
			}
			log.SetOutput(os.Stdout)
			time.Sleep(cfg.Interval)
			continue
		}
		if !cfg.Quiet {
			log.Printf("Probe #%d (%d success, %d fail): %s", probeCount, probeCount-failCount, failCount, fd)
		}
		time.Sleep(cfg.Interval)
	}
	return failCount, nil
}
