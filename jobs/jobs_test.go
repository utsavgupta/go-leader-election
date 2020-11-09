package jobs_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/utsavgupta/go-leader-election/jobs"
)

func TestNewPreacher(t *testing.T) {

	t.Parallel()

	message := "this is preachy"

	var tests = []struct {
		times          int
		expectedOutput string
	}{
		{times: 0, expectedOutput: ""},
		{times: 1, expectedOutput: "this is preachy\n"},
		{times: 2, expectedOutput: "this is preachy\nthis is preachy\n"},
	}

	for _, test := range tests {
		var b []byte
		bufW := bytes.NewBuffer(b)

		preacher := jobs.NewPreacher("test preacher", message, bufW)

		preacher.Do(context.Background(), test.times)

		if actualOutput := bufW.String(); test.expectedOutput != actualOutput {
			t.Fatalf("Expected output: %s, Actual output: %s\n", test.expectedOutput, actualOutput)
		}
	}
}
