package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiTarget = "mf"
	apiQuery  = `{ "queryDynamicShadowData" : 1 }`
)

// Config holds the configuration parameters used throughout the prober.
type Config struct {
	Host        string
	Interval    time.Duration
	ExitOnError bool
	Quiet       bool
	Timeout     time.Duration
}

func (cfg Config) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("Target: %q; ", cfg.Host))
	s.WriteString(fmt.Sprintf("Interval: %s; ", cfg.Interval))
	s.WriteString(fmt.Sprintf("Timeout: %s; ", cfg.Timeout))
	s.WriteString(fmt.Sprintf("Exit on error: %v; ", cfg.ExitOnError))
	s.WriteString(fmt.Sprintf("Quiet mode: %v", cfg.Quiet))

	return s.String()
}

// FanData holds the basic information returned by the fan.
type FanData struct {
	ClientID        string `json:"clientID"`
	FanOn           bool   `json:"fanOn"`
	FanSpeed        int    `json:"fanSpeed"`
	FanDirection    string `json:"fanDirection"`
	Wind            bool   `json:"wind"`
	WindSpeed       int    `json:"windSpeed"`
	LightOn         bool   `json:"lightOn"`
	LightBrightness int    `json:"lightBrightness"`
}

func (fd FanData) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("ClientID: %q; ", fd.ClientID))
	s.WriteString(fmt.Sprintf("Fan: %v (speed: %d, direction: %q); ", fd.FanOn, fd.FanSpeed, fd.FanDirection))
	s.WriteString(fmt.Sprintf("Wind: %v (speed: %d); ", fd.Wind, fd.WindSpeed))
	s.WriteString(fmt.Sprintf("Light: %v (brightness: %d)", fd.LightOn, fd.LightBrightness))

	return s.String()
}

var cfg Config

func init() {
	var iv, to int

	flag.StringVar(&cfg.Host, "host", "", "Host to probe")
	flag.IntVar(&iv, "interval", 10, "Polling interval in seconds")
	flag.BoolVar(&cfg.ExitOnError, "exit-on-error", true, "Exit polling loop on error")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "Log only polling errors")
	flag.IntVar(&to, "timeout", 1, "Polling probe timeout in seconds")
	flag.Parse()

	if cfg.Host == "" {
		log.Fatalf("Target host must be specified.")
	}
	cfg.Interval = time.Duration(iv) * time.Second
	cfg.Timeout = time.Duration(to) * time.Second
	log.SetOutput(os.Stdout)
}

// PollFan executes HTTP POST to query the fan status, reporting error.
func PollFan(host string) (*FanData, error) {
	fd := new(FanData)

	url := fmt.Sprintf("http://%s/%s", cfg.Host, apiTarget)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(apiQuery)))
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

func main() {
	if !cfg.Quiet {
		log.Print(cfg)
	}

	probeCount := 1
	for {
		fd, err := PollFan(cfg.Host)
		switch err != nil {
		case true:
			{
				log.SetOutput(os.Stderr)
				log.Printf("Probe #%d: error reading fan information from %q: %v.", probeCount, cfg.Host, err)
				if cfg.ExitOnError {
					os.Exit(42)
				}
				log.SetOutput(os.Stdout)
			}
		case false:
			{
				if !cfg.Quiet {
					log.Printf("Probe #%d: %s", probeCount, fd)
				}
			}
			probeCount++
			time.Sleep(cfg.Interval)
		}
	}
}
