package schedulers

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/utsavgupta/go-leader-election/globals"
	"github.com/utsavgupta/go-leader-election/jobs"
)

// Scheduler type function consists of a context c, a job j, and a duration t.
// Scheduler functions perform j every t time units until c is cancelled.
type Scheduler func(context.Context, jobs.Job, time.Duration)

type lease struct {
	Leader string
	Expiry time.Time
}

// NewScheduler returns a new scheduler with the given name. It uses the
// given client instance for contesting leadership for a job.
func NewScheduler(nodeName string, client *datastore.Client) Scheduler {
	return func(ctx context.Context, job jobs.Job, t time.Duration) {
		globals.Logger.Infof(ctx, "Starting %s scheduler on node %s", job.Name, nodeName)
		for {
			select {
			case <-time.After(t):
				err := becomeLeader(ctx, nodeName, job.Name, client, t)

				if err == nil {
					job.Do(ctx, 5)
				} else {
					globals.Logger.Errorf(ctx, "Idle: %s", err)
				}
			case <-ctx.Done():
				globals.Logger.Infof(ctx, "Terminating %s scheduler on node %s", job.Name, nodeName)
				return
			}
		}
	}
}

func becomeLeader(ctx context.Context, nodeName, jobName string, client *datastore.Client, t time.Duration) error {

	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var l lease

		key := datastore.NameKey("Lease", jobName, nil)

		err := tx.Get(key, &l)

		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		// Become the leader only if an entry for the given job does not exist OR the lease of the previous
		// leader has already expired OR the current scheduler was the previous leader
		if err == datastore.ErrNoSuchEntity || l.Expiry.Before(time.Now()) || l.Leader == nodeName {
			l.Leader = nodeName
			l.Expiry = time.Now().Add(t)
			_, err := tx.Put(key, &l)

			return err
		}

		return fmt.Errorf("Node %s could not become leader for job %s", nodeName, jobName)
	})

	return err
}
