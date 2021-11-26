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
	cfg                mffp.Config
	version, gitCommit string
)

func printVersion(v, g string) {
	fmt.Printf("Version: %s\n Commit: %s\n", v, g)
}

func init() {
	var v bool

	flag.StringVar(&cfg.Host, "host", "", "Host to probe")
	flag.DurationVar(&cfg.Interval, "interval", 10*time.Second, "Polling interval in seconds")
	flag.BoolVar(&cfg.ExitOnError, "exit-on-error", true, "Exit polling loop on error")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "Log only polling errors")
	flag.DurationVar(&cfg.Timeout, "timeout", 1*time.Second, "Polling probe timeout in seconds")
	flag.IntVar(&cfg.MaxCount, "maxcount", 0, "Maximum number of probes (0 means unlimited)")
	flag.IntVar(&cfg.MaxFail, "maxfail", 0, "Maximum number of failed probes")
	flag.BoolVar(&v, "version", false, "Show version")
	flag.Parse()

	if v {
		printVersion(version, gitCommit)
		os.Exit(42)
	}

	if cfg.MaxCount < 0 {
		log.Fatalf("Maximum probe count must be greater or equal to 0.")
	}

	if cfg.MaxFail < 0 {
		log.Fatalf("Maximum probe number of failed probes be greater or equal to 0.")
	}

	if cfg.Host == "" {
		log.Fatalf("Target host must be specified.")
	}
	log.SetOutput(os.Stdout)
}

func main() {
	if !cfg.Quiet {
		log.Print(cfg)
	}

	cmd.ProbeFan(&cfg)
}
