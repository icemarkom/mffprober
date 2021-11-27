package mffprober

import (
	"fmt"
	"strings"
	"time"
)

const (
	// MFFTarget is the API target in URL.
	MFFTarget = "mf"
	// MFFQuery is the API call to query the fan status.
	MFFQuery = `{ "queryDynamicShadowData" : 1 }`
	// MFFReboot is the API call to reboot the fan.
	MFFReboot = `{ "reboot" : true }`
)

// Config holds the configuration parameters used throughout the prober.
type Config struct {
	Host                string
	Interval            time.Duration
	Quiet               bool
	Timeout             time.Duration
	Count, MaxFailCount int
	Reboot              bool
}

func (cfg Config) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("Target: %q; ", cfg.Host))
	s.WriteString(fmt.Sprintf("Interval: %s; ", cfg.Interval))
	s.WriteString(fmt.Sprintf("Timeout: %s; ", cfg.Timeout))
	s.WriteString(fmt.Sprintf("Quiet mode: %v; ", cfg.Quiet))
	s.WriteString(fmt.Sprintf("Probe count: %v; ", cfg.Count))
	s.WriteString(fmt.Sprintf("Max failed probes: %v", cfg.MaxFailCount))

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
