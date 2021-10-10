package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func Client(cfg aws.Config) *iam.Client {
	return iam.NewFromConfig(cfg)
}

type new_key struct {
	key    string
	secret string
}

func CreateKey(ctx context.Context, client *iam.Client) new_key {
	keys, err := client.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{})

	if err != nil {
		Boom("Can't create new key", err)
	}

	return new_key{
		key:    *keys.AccessKey.AccessKeyId,
		secret: *keys.AccessKey.SecretAccessKey,
	}
}

func DeleteKey(ctx context.Context, client *iam.Client, key_id string) error {
	input := &iam.DeleteAccessKeyInput{AccessKeyId: &key_id}
	_, err := client.DeleteAccessKey(ctx, input)

	return err
}
