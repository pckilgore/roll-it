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

func EnvironmentCredentials(ctx context.Context) Credentials {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		Boom("There was an error loading configuration:", err)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		Boom("There was an error loading credentials chain:", err)
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
		Boom(msg, err)
	}
	return CredentialsFile{path, string(cfg)}
}

func WriteCredentialsFile(to CredentialsFile) {
	err := ioutil.WriteFile(to.filepath, []byte(to.content), 644)

	if err != nil {
		note := "------Unable to write credentials: Recovery Information------"
		fmt.Printf("%s\n%s\n%s\n\n", note, to.content, note)
		Boom("Failed to write new credentials to file.", err)
	}
}
