package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"io/ioutil"
	"regexp"
	"strings"
)

type keyLocator struct {
	key string
	re  *regexp.Regexp
}

var AccessKeyLocator = &keyLocator{
	"aws_access_key_id",
	regexp.MustCompile("(?P<key>aws_access_key_id\\s{0,}=\\s{0,})(?P<value>[A-Z0-9]{20})"),
}

var SecretKeyLocator = &keyLocator{
	"aws_secret_access_key",
	regexp.MustCompile("(?P<key>aws_secret_access_key\\s{0,}=\\s{0,})(?P<value>[A-Za-z0-9/+=]{40})"),
}

func (loc *keyLocator) Replace(replace, with, in string) string {
	return loc.re.ReplaceAllStringFunc(in, func(match string) string {
		if strings.Contains(match, replace) {
			return fmt.Sprintf("%s=%s", loc.key, with)
		}
		return match
	})
}

func (loc *keyLocator) MatchKey(find string, in string) (match string) {
	matched_keys := loc.re.FindAllStringSubmatch(in, -1)
	key_value_group := loc.re.SubexpIndex("value")
	for _, groups := range matched_keys {
		val := groups[key_value_group]
		if val == find {
			match = val
		}
	}
	return
}

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Config          aws.Config
}

func EnvironmentCredentials() Credentials {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		Boom(err, "There was an error loading configuration:")
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		Boom(err, "There was an error loading credentials chain:")
	}

	return Credentials{creds.AccessKeyID, creds.SecretAccessKey, cfg}
}

type CredentialsFile struct {
	filepath string
	content  string
}

func GetCredentialsFile() CredentialsFile {
	path := config.DefaultSharedCredentialsFilename()
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf("Can't load config from %s", path)
		Boom(err, msg)
	}
	return CredentialsFile{path, string(cfg)}
}
