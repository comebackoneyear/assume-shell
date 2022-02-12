package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/comebackoneyear/assume-shell/pkg/profile"
	"github.com/comebackoneyear/assume-shell/pkg/shell"
)

var (
	// Populated by goreleaser during build
	version = "dev"
	commit  = ""
	date    = ""
	arch    = ""
)

func main() {
	flagVersion := flag.Bool("v", false, "prints current assume-shell version")

	flag.Parse()

	if *flagVersion {
		if commit != "" {
			fmt.Printf("assume-shell %s arch:%s commit:%s date:%s\n", version, arch, commit, date)
		} else {
			fmt.Printf("assume-shell %s build\n", version)
		}
		os.Exit(0)
	}

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
	fmt.Fprintf(os.Stderr, "Usage: %s [-v|<profile name>]\n", os.Args[0])
	flag.PrintDefaults()
}
