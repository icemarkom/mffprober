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
	var v, e bool

	flag.StringVar(&cfg.Host, "host", "", "Host to probe")
	flag.DurationVar(&cfg.Interval, "interval", 10*time.Second, "Polling interval in seconds")
	flag.BoolVar(&e, "exit-on-error", true, "DEPRECATED: Use maxfail=0 to disable, or a positive value to control")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "Log only polling errors")
	flag.DurationVar(&cfg.Timeout, "timeout", 1*time.Second, "Polling probe timeout in seconds")
	flag.IntVar(&cfg.Count, "count", 0, "Maximum number of probes (0 means unlimited)")
	flag.IntVar(&cfg.MaxFailCount, "maxfail", 1, "Maximum number of failed probes (0 means unlimited)")
	flag.BoolVar(&v, "version", false, "Show version")
	flag.Parse()

	if v {
		printVersion(version, gitCommit)
		os.Exit(42)
	}

	if !e {
		log.SetOutput(os.Stderr)
		log.Printf("You have used deprecated %q flag. Use %q in the future.", "exit-on-error", "maxfail=0")
		cfg.MaxFailCount = 0
		log.SetOutput(os.Stderr)
	}

	if cfg.Count < 0 {
		log.Fatalf("Probe count must be greater or equal to 0.")
	}

	if cfg.MaxFailCount < 0 {
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
	os.Exit(cmd.ProbeFan(&cfg))
}
