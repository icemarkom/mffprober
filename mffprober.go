// Copyright 2021 MFFProber Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

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
	if !cfg.Reboot {
		s.WriteString(fmt.Sprintf("Probe count: %v; ", cfg.Count))
		s.WriteString(fmt.Sprintf("Max failed probes: %v", cfg.MaxFailCount))
	}
	s.WriteString(fmt.Sprintf("Reboot: %v", cfg.Reboot))

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
