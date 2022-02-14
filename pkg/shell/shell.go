package shell

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hako/durafmt"
)

func ShellWithCredentials(profile string, creds aws.Credentials) error {
	shell := os.Getenv("SHELL")

	argv0, err := exec.LookPath(shell)
	if err != nil {
		return err
	}

	env := os.Environ()

	env = append(env, []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken),
		fmt.Sprintf("AWS_SECURITY_TOKEN=%s", creds.SessionToken),
		fmt.Sprintf("ASSUMED_PROFILE=%s", profile),
	}...)

	var expires string
	if creds.CanExpire {
		env = append(env, fmt.Sprintf("ASSUMED_PROFILE_EXPIRES=%d", creds.Expires.Unix()))
		expires = fmt.Sprintf("expires in %s", durafmt.Parse(time.Until(creds.Expires)).LimitFirstN(2))
	} else {
		expires = "never expires"
	}

	fmt.Printf("Exported assumed role credentials for profile %s, %s\n", profile, expires)

	cmd := exec.Command(argv0)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	inputSignal := make(chan os.Signal, 1)
	signal.Notify(inputSignal, os.Interrupt, syscall.SIGTERM)

	if err := cmd.Start(); err != nil {
		return err
	}

	waitChannel := make(chan error, 1)
	go func() {
		waitChannel <- cmd.Wait()
		close(waitChannel)
	}()

	for {

		select {
		case sig := <-inputSignal:
			if err := cmd.Process.Signal(sig); err != nil {
				return err
			}
		case err := <-waitChannel:
			var waitStatus syscall.WaitStatus
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus = exitError.Sys().(syscall.WaitStatus)
				os.Exit(waitStatus.ExitStatus())
			}
			if err != nil {
				return err
			}
			return nil
		}
	}

}
