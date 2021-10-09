package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Timeout SDK interactions in go routines after ten seconds.
	// This should take less than a second, so something is horribly wrong.
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Parallelize credential chain + credentials file load operations
	creds_file_load := make(chan CredentialsFile)
	cred_load := make(chan Credentials)

	go func() {
		creds_file_load <- GetCredentialsFile()
	}()
	go func() {
		cred_load <- EnvironmentCredentials(ctx)
	}()

	creds_file := <-creds_file_load
	creds := <-cred_load

	client := Client(creds.Config)

	// Identify correct credentials to rotate
	key_match := AccessKeyLocator.MatchKey(creds.AccessKeyID, creds_file.content)
	secret_match := SecretKeyLocator.MatchKey(creds.SecretAccessKey, creds_file.content)

	if len(secret_match) == 0 || len(key_match) == 0 {
		tmpl := "Cannot find key_id=%s in file=%s)"
		msg := fmt.Sprintf(tmpl, creds.AccessKeyID, creds_file.filepath)
		Boom("Aborting Intentionally:", fmt.Errorf(msg))
	}

	// Create new Credentials
	new_key := CreateKey(ctx, client)

	// Parallelize deletion of current credentials
	del_op := make(chan bool)
	go func() {
		DeleteKey(ctx, client, string(key_match))
		del_op <- true
	}()

	// Write out new credentials
	creds_file.content = AccessKeyLocator.Replace(key_match, new_key.key, creds_file.content)
	creds_file.content = SecretKeyLocator.Replace(secret_match, new_key.secret, creds_file.content)
	WriteCredentialsFile(creds_file)

	// Wait for delete routine to complete
	<-del_op
}
