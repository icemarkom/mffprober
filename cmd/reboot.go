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

package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	mffp "github.com/icemarkom/mffprober"
)

// // RebootFan executes HTTP POST to reboot the fan.
func RebootFan(cfg *mffp.Config) (int, error) {

	u := fmt.Sprintf("http://%s/%s", cfg.Host, mffp.MFFTarget)

	req, err := http.NewRequest("POST", u, bytes.NewBuffer([]byte(mffp.MFFReboot)))
	if err != nil {
		return 1, fmt.Errorf("error formatting request: %w", err)
	}

	hc := &http.Client{
		Timeout: cfg.Timeout,
	}
	r, err := hc.Do(req)
	if err != nil {
		// Fan does not respond to a reboot request, it just times-out.
		// We treat that as a potential success. Otherwise, it's an error.
		if err, ok := err.(*url.Error); ok && !err.Timeout() {
			return 1, fmt.Errorf("error rebooting the fan: %w", err)
		}
		r = new(http.Response)
		r.StatusCode = http.StatusRequestTimeout
	}
	// Just in a case the fans get API fix, treat StatusOK andd StatusRequestTimeout as success.
	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusRequestTimeout {
		return 1, fmt.Errorf("fan reported HTTP error: %s", r.Status)
	}
	if !cfg.Quiet {
		log.Printf("Reboot request to %q sent.", cfg.Host)
	}
	time.Sleep(cfg.Interval)
	fd, err := PollFan(cfg)
	if err != nil {
		return 1, fmt.Errorf("probably failed to reboot the fan: %w", err)
	}
	if !cfg.Quiet {
		log.Printf("Fan is alive: %s", fd)
	}
	return 0, nil
}
