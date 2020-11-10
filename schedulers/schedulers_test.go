package schedulers_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/utsavgupta/go-leader-election/globals"

	"cloud.google.com/go/datastore"
	"github.com/utsavgupta/go-leader-election/jobs"
	"github.com/utsavgupta/go-leader-election/schedulers"
)

const (
	PROJECT_ID         = "go-leader-election"
	EMULATOR_HOST      = "localhost:8081"
	EMULATOR_HOST_PATH = "localhost:8081/datastore"
	DATASTORE_HOST     = "http://localhost:8081"
)

func TestNewScheduler(t *testing.T) {

	t.Parallel()

	ctx, fnCancel := context.WithCancel(context.Background())
	defer fnCancel()

	os.Setenv("DATASTORE_DATASET", PROJECT_ID)
	os.Setenv("DATASTORE_EMULATOR_HOST", EMULATOR_HOST)
	os.Setenv("DATASTORE_EMULATOR_HOST_PATH", EMULATOR_HOST_PATH)
	os.Setenv("DATASTORE_HOST", DATASTORE_HOST)
	os.Setenv("DATASTORE_PROJECT_ID", PROJECT_ID)

	globals.InitLogger("scheduler-test", "test")

	client, err := datastore.NewClient(ctx, "")

	if err != nil {
		panic(err)
	}

	var k int
	mutex := &sync.Mutex{}

	job := jobs.Job{
		Name: "incrementer",
		Do: func(ctx context.Context, n int) {
			defer mutex.Unlock()
			mutex.Lock()

			k += 1
		},
	}

	time.AfterFunc(15*time.Second, func() {
		fnCancel()
	})

	// spin of then schedulers for the same job
	for i := 1; i <= 10; i++ {
		func(n int) {
			s := schedulers.NewScheduler(fmt.Sprintf("scheduler-%d", i), client)

			go s(ctx, job, 2*time.Second)
		}(i)
	}

	time.Sleep(15 * time.Second)

	// The jobs are scheduled every 2 seconds. We wait for a total of 15 seconds.
	// Hence the total no of jobs that should get scheduled should be floor(15/2)
	// or 7
	if k != 7 {
		t.Fatalf("Expected 7 job executions, got %d", k)
	}
}
