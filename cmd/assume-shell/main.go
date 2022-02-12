package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/comebackoneyear/assume-shell/pkg/profile"
	"github.com/comebackoneyear/assume-shell/pkg/shell"
)

func main() {
	flag.Parse()
	argv := flag.Args()

	assumeProfile := "default"

	if len(argv) == 1 {
		assumeProfile = argv[0]

	}

	if len(argv) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	creds, err := profile.AssumeProfile(assumeProfile)
	if err != nil {
		log.Fatal(err)
	}

	err = shell.ShellWithCredentials(assumeProfile, creds)
	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [profile]\n", os.Args[0])
	flag.PrintDefaults()
}
