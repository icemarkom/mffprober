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
}

func (cfg Config) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("Target: %q; ", cfg.Host))
	s.WriteString(fmt.Sprintf("Interval: %s; ", cfg.Interval))
	s.WriteString(fmt.Sprintf("Exit on error: %v", cfg.ExitOnError))

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
	var iv int

	flag.StringVar(&cfg.Host, "host", "", "Host to probe")
	flag.IntVar(&iv, "interval", 10, "Polling interval in seconds")
	flag.BoolVar(&cfg.ExitOnError, "exit-on-error", true, "Exit polling loop on error")
	flag.Parse()

	if cfg.Host == "" {
		log.Fatalf("Target host must be specified.")
	}

	cfg.Interval = time.Duration(iv) * time.Second
}

// PollFan executes HTTP POST to query the fan status, reporting error.
func PollFan(host string) (*FanData, error) {
	fd := new(FanData)

	url := fmt.Sprintf("http://%s/%s", cfg.Host, apiTarget)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(apiQuery)))
	if err != nil {
		return nil, fmt.Errorf("error formatting request: %w", err)
	}

	hc := &http.Client{}
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
	log.Print(cfg)

	for {
		fd, err := PollFan(cfg.Host)
		switch err != nil {
		case true:
			{
				log.Printf("Error reading fan information from %q: %v.", cfg.Host, err)
				if cfg.ExitOnError {
					os.Exit(42)
				}
			}
		case false:
			{
				log.Print(fd)
			}
			time.Sleep(cfg.Interval)
		}
	}
}
