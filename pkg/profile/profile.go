package profile

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func AssumeProfile(profile string) (aws.Credentials, error) {
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
