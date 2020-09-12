package main

import (
	"flag"
	"log"
	"os"
	"time"

	mffp "github.com/icemarkom/mffprober"
	"github.com/icemarkom/mffprober/cmd"
)

var cfg mffp.Config

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

func main() {
	if !cfg.Quiet {
		log.Print(cfg)
	}

	cmd.ProbeFan(&cfg)
}
