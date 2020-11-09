package jobs

import (
	"context"
	"io"
	"time"

	"github.com/utsavgupta/go-leader-election/globals"
)

// Job type contains a name and a doable function that is repeated N times
type Job struct {
	Name string
	Do   func(context.Context, int)
}

// NewPreacher returns a job that preaches the given message
// by writing it to a string writer N times. Sleeping for exactly
// one second between each operation.
func NewPreacher(name, message string, writer io.StringWriter) Job {
	return Job{
		Name: name,
		Do: func(ctx context.Context, times int) {

			// Don't sleep before writing the first message
			sleep := false

			for i := 1; i <= times; i++ {

				if sleep {
					time.Sleep(1 * time.Second)
				}

				_, err := writer.WriteString(message + "\n")

				if err != nil {
					globals.Logger.Error(ctx, err)
				}

				sleep = true
			}
		},
	}
}
