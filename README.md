# assume-shell
This tool will request AWS credentials for a given profile/role and start a new shell with them in the environment.

## Installation

If you have a working Go >=1.17 environment:

```bash
$ go install -u github.com/comebackoneyear/assume-shell/cmd/assume-shell@latest
```

If you have a working Go <1.17 environment:

```bash
$ go get -u github.com/comebackoneyear/assume-shell/cmd/assume-shell
```
## Configuration

Setup a profile for each role you would like to assume in `~/.aws/config`.

For example:

`~/.aws/config`:

```ini
[profile work]
region = eu-north-1

[profile stage]
# Stage AWS Account.
region = eu-west-1
role_arn = arn:aws:iam::00000001234:role/Admin
source_profile = work

[profile prod]
# Production AWS Account.
region = us-east-1
role_arn = arn:aws:iam::00000005678:role/Deploy
source_profile = work
```

`~/.aws/credentials`:

```ini
[work]
aws_access_key_id = <AWS_ACCESS_KEY_ID>
aws_secret_access_key = <AWS_SECRET_ACCESS_KEY>
```

Reference: https://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html

In this example, we have three AWS Account profiles:

 * work
 * stage
 * prod

Each member of the org has their own IAM user and access/secret key for the `work` AWS Account.
The keys are stored in the `~/.aws/credentials` file.

The `stage` and `prod` AWS Accounts have IAM roles named `Admin` and `Deploy`.
The `assume-shell` tool helps a user authenticate (using their keys) and then assume the privilege of the the role, even across AWS accounts!

## Usage

Start a new shell with the role stage:

```bash
$ assume-shell stage
```

The `assume-shell` tool sets `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN` environment variables and then executes the $SHELL. In addition the `ASSUMED_PROFILE` variable will be set to whatever profile was assumed

## TODO

* [ ] Use default profile on empty argv
* [ ] Add MFA support
* [ ] Add brew installer or MacOS
* [ ] Add Support to execute commands
* [ ] Add support to configure shell prompt to show active profile
* [ ] Test and add support for multiple shells (zsh, fish, powershell)


## Credits
Heavily inspired by https://github.com/remind101/assume-role
