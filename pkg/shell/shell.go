package shell

import (
	"fmt"
	"os"
	"os/exec"
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

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
	os.Setenv("AWS_SECURITY_TOKEN", creds.SessionToken)
	os.Setenv("ASSUMED_PROFILE", profile)

	var expires string
	if creds.CanExpire {
		os.Setenv("ASSUMED_PROFILE_EXPIRES", fmt.Sprintf("%d", creds.Expires.Unix()))
		expires = fmt.Sprintf("expires in %s", durafmt.Parse(time.Until(creds.Expires)).LimitFirstN(2))
	} else {
		expires = "never expires"
	}
	env := os.Environ()
	fmt.Printf("Exported assumed role credentials for profile %s, %s\n", profile, expires)
	return syscall.Exec(argv0, []string{}, env)
}
