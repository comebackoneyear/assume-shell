package shell

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func ShellWithCredentials(profile string, creds aws.Credentials) error {
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
