package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"os"
)

func Client(cfg aws.Config) *iam.Client {
	return iam.NewFromConfig(cfg)
}

func ListKeys(client *iam.Client) chan []types.AccessKeyMetadata {
	r := make(chan []types.AccessKeyMetadata)

	go func() {
		input := &iam.ListAccessKeysInput{}
		keys, err := client.ListAccessKeys(context.TODO(), input)

		if err != nil {
			panic("Can't find keys")
		}
		r <- keys.AccessKeyMetadata
	}()

	return r
}

type new_key struct {
	key    string
	secret string
}

func CreateKey(client *iam.Client) new_key {
	keys, err := client.CreateAccessKey(context.TODO(), &iam.CreateAccessKeyInput{})

	if err != nil {
		Boom(err, "Can't create new key")
	}

	return new_key{
		key:    *keys.AccessKey.AccessKeyId,
		secret: *keys.AccessKey.SecretAccessKey,
	}
}

func DeleteKey(client *iam.Client, key_id string) {
	input := &iam.DeleteAccessKeyInput{AccessKeyId: &key_id}
	_, err := client.DeleteAccessKey(context.TODO(), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not delete old key: %v\n", err)
		os.Exit(1)
	}
}
