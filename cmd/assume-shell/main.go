package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	flag.Parse()
	argv := flag.Args()

	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	profile := argv[0]
	creds, err := assumeProfile(profile)
	if err != nil {
		log.Fatal(err)
	}

	err = shellWithCredentials(profile, creds)
	if err != nil {
		log.Fatal(err)
	}

}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <profile>\n", os.Args[0])
	flag.PrintDefaults()
}

func assumeProfile(profile string) (aws.Credentials, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return aws.Credentials{}, err
	}

	return cfg.Credentials.Retrieve(ctx)
}

func shellWithCredentials(profile string, creds aws.Credentials) error {
	shell := os.Getenv("SHELL")

	argv0, err := exec.LookPath(shell)
	if err != nil {
		return err
	}

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
	os.Setenv("AWS_SECURITY_TOKEN", creds.SessionToken)
	os.Setenv("ASSUMED_PROFILE", profile)

	env := os.Environ()
	return syscall.Exec(argv0, []string{}, env)
}
