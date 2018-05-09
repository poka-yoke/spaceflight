package http

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
)

// Backoff holds configuration for the backoff mechanism for http requests
type Backoff struct {
	ready         func() bool   // Function to determine when we are done
	wait          time.Duration // How much time to wait between retries
	timelimit     time.Duration // Maximum time to wait between retries
	withTimeLimit bool          // Whether we have a time limit
}

// NewConstantBackoff returns a new backoff with constant waiting time
func NewConstantBackoff(ready func() bool, sleep time.Duration) *Backoff {
	return &Backoff{
		ready:         ready,
		wait:          sleep,
		withTimeLimit: false,
	}
}

// NewLinearBackoff returns a new backoff with linearly increasing waiting time
// The waiting time is increased by the base `sleep` amount after
// every iteration, unless the `limit` has already been reached.  This
// effectively means that the function can sleep for longer than
// `limit` but never for as long as `limit + sleep`.
func NewLinearBackoff(ready func() bool, sleep, timelimit time.Duration) *Backoff {
	return &Backoff{
		ready:         ready,
		wait:          sleep,
		timelimit:     timelimit,
		withTimeLimit: true,
	}
}

// Do executes the request with backoff enabled
func (b Backoff) Do(req *request.Request) {
	sleep := b.wait
	for {
		err := req.Send()
		if err != nil {
			continue
		}
		if b.ready() {
			return
		}
		log.Printf("Waiting for %fs\n", sleep.Seconds())
		time.Sleep(sleep)
		sleep = b.waitDuration(sleep)
	}
}

// waitDuration tells us how long do we need to wait
func (b Backoff) waitDuration(sleep time.Duration) time.Duration {
	if !b.withTimeLimit {
		return b.wait
	}
	if sleep < b.timelimit {
		return sleep + b.wait
	}
	return sleep
}
