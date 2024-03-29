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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	mffp "github.com/icemarkom/mffprober"
	"github.com/icemarkom/mffprober/cmd"
)

var (
	cfg        mffp.Config
	version    = "development"
	gitCommit  = ""
	binaryName = "mffprober"
)

func printVersion() {
	fmt.Fprint(flag.CommandLine.Output(), programVersion())
}

func init() {
	var (
		e, v bool
		h    string
	)
	flag.Usage = func() { programUsage() }

	flag.StringVar(&h, "host", "", "DEPRECATED: Specify host name/address after the flags")
	flag.DurationVar(&cfg.Interval, "interval", 10*time.Second, "Polling interval in seconds")
	flag.BoolVar(&e, "exit-on-error", true, "DEPRECATED: Use maxfail=0 to disable, or a positive value to control")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "Log only polling errors")
	flag.DurationVar(&cfg.Timeout, "timeout", 1*time.Second, "Polling probe timeout in seconds")
	flag.IntVar(&cfg.Count, "count", 0, "Maximum number of probes (0 means unlimited)")
	flag.IntVar(&cfg.MaxFailCount, "maxfail", 1, "Maximum number of failed probes (0 means unlimited)")
	flag.BoolVar(&cfg.Reboot, "reboot", false, "Reboot fan. Ignores most flags")
	flag.BoolVar(&v, "version", false, "Show version")
	flag.Parse()

	if v {
		printVersion()
		os.Exit(42)
	}

	// Handle deprecated "host" flag gracefully.
	if h != "" {
		log.SetOutput(os.Stderr)
		log.Printf("You have used deprecated %q flag. Specify host name/address after the flags.", "host")
		cfg.Host = h
	}

	// Handle deprecated "exit-on-error" flag gracefully.
	if !e {
		log.SetOutput(os.Stderr)
		log.Printf("You have used deprecated %q flag. Use %q in the future.", "exit-on-error", "maxfail=0")
		cfg.MaxFailCount = 0
	}

	if cfg.Count < 0 {
		log.Fatalf("Probe count must be greater or equal to 0.")
	}

	if cfg.MaxFailCount < 0 {
		log.Fatalf("Maximum probe number of failed probes be greater or equal to 0.")
	}

	// More complex logic here needed to allow for using --host.
	if cfg.Host == "" {
		if flag.NArg() == 0 {
			log.SetOutput(os.Stderr)
			log.Fatalf("Target host must be specified.")
		}
		cfg.Host = flag.Arg(0)
	}

	log.SetOutput(os.Stdout)
}

func main() {
	if !cfg.Quiet {
		log.Print(cfg)
	}
	command := cmd.ProbeFan
	if cfg.Reboot {
		command = cmd.RebootFan
	}
	fc, err := command(&cfg)
	if err != nil {
		if !cfg.Quiet {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	if fc > 255 {
		os.Exit(255)
	}
	os.Exit(fc)
}
