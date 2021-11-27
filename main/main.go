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

func printVersion(v, g string) {
	fmt.Printf("Version: %s\n Commit: %s\n", v, g)
}

func init() {
	var (
		v, e bool
		h    string
	)
	flag.Usage = func() { programUsage() }

	flag.StringVar(&h, "host", "", "DEPRECATED: Specify host name/address after the flags.")
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
	fc := cmd.ProbeFan(&cfg)
	if fc > 255 {
		os.Exit(255)
	}
	os.Exit(fc)
}
