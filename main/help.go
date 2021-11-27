package main

import (
	"flag"
	"fmt"
	"strings"
)

func programVersion() string {
	return fmt.Sprintf("Version: %s\n Commit: %s\n", version, gitCommit)
}

func programUsage() {
	s := new(strings.Builder)

	s.WriteString(fmt.Sprintf("Usage: %s [flags] hostname.or.ip.address\n\nFlags:\n", binaryName))

	o := flag.CommandLine.Output()
	flag.CommandLine.SetOutput(s)
	flag.PrintDefaults()
	flag.CommandLine.SetOutput(o)

	s.WriteString(fmt.Sprintln())
	s.WriteString(programVersion())

	fmt.Fprint(flag.CommandLine.Output(), s.String())
}
