package tools

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFallbackLogger(t *testing.T) {
	type loggerTest struct {
		msg    string
		output string
	}

	loggerTests := []loggerTest{
		{"Sample Error", " level=error msg=\"Sample Error\"\n"},
		{"Another Error", " level=error msg=\"Another Error\"\n"},
	}

	for _, lt := range loggerTests {
		rescueStdout := os.Stderr
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("failed to get pipe: %v", err)
		}

		os.Stderr = w

		FallbackLogger(errors.New(lt.msg))

		w.Close()

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read data: %v", err)
		}

		os.Stderr = rescueStdout

		s := string(out)
		if !strings.HasSuffix(string(s), lt.output) {
			t.Fatalf("expected %s, got %s", lt.output, s)
		}
	}
}
