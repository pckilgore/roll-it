package main

import (
	"fmt"
	"os"
)

func Boom(msg string, err error) {
	fmt.Fprintf(
		os.Stderr,
		"\n%s\n%+v\n\n%s\n",
		msg,
		err,
		"Stopping for manual intervention",
	)
	os.Exit(1)
}
