// Package tools contains useful utils to reduce amount of boilerplate code.
package tools

import (
	"fmt"
	"os"
	"time"
)

// FallbackLogger is emergency logger for case when init of normal logger was failed.
func FallbackLogger(msg error) {
	ts := time.Now().UTC().Format("2006-01-02T15:04:05-0700")

	fmt.Fprintf(os.Stderr, "time=%q level=error msg=%q\n", ts, msg)
}
