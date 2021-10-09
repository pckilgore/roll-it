package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Boom(err error, msg_op ...string) {
	msg := "There was an error:"
	if len(msg_op) > 0 {
		msg = msg_op[0]
	}

	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "%v\n", err)
	fmt.Fprintln(os.Stderr, "Stopping for manual intervention")
	os.Exit(1)
}

func main() {
	creds_file_load := make(chan CredentialsFile)
	cred_load := make(chan Credentials)

	go func() {
		creds_file_load <- GetCredentialsFile()
	}()
	go func() {
		cred_load <- EnvironmentCredentials()
	}()

	creds_file := <-creds_file_load
	creds := <-cred_load
	client := Client(creds.Config)

	key_match := AccessKeyLocator.MatchKey(creds.AccessKeyID, creds_file.content)
	secret_match := SecretKeyLocator.MatchKey(creds.SecretAccessKey, creds_file.content)

	if len(secret_match) == 0 || len(key_match) == 0 {
		msg := `Could not match credentials loaded in chain to credentials file: %s`
		Boom(fmt.Errorf(msg, creds_file.filepath))
	}

	new_key := CreateKey(client)

	del_op := make(chan bool)
	go func() {
		DeleteKey(client, string(key_match))
		del_op <- true
	}()

	creds_file.content = AccessKeyLocator.Replace(key_match, new_key.key, creds_file.content)
	creds_file.content = SecretKeyLocator.Replace(secret_match, new_key.secret, creds_file.content)

	err := ioutil.WriteFile(creds_file.filepath, []byte(creds_file.content), 644)

	if err != nil {
		fmt.Println("Failed to write new credentials to file. Investigate manually.")
	}

	<-del_op
}
