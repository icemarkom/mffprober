package main

import (
	"flag"
	"fmt"
)

func programVersion() {
	fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n Commit: %s\n", version, gitCommit)
}

func programUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] hostname.or.ip.address\n\nFlags:\n", binaryName)
	flag.PrintDefaults()
	fmt.Fprintln(flag.CommandLine.Output())
	programVersion()
}
