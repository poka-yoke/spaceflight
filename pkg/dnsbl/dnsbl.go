package dnsbl

import (
	"net"
	"sync"
)

// Checker controls the flow of package and provides a single point of
// entry for its users.
type Checker struct {
	// length is the number of providers
	// queried is the number of providers who answered
	// positive is the number of appearances reported
	length, queried, positive int
	providers                 []string

	// Functional control
	wg sync.WaitGroup
}

// NewChecker creates a new, default configured Checker
func NewChecker(providers []string) *Checker {
	return &Checker{providers: providers}
}

// Query contacts the providers to check if the IP is present in their
// lists
func (c *Checker) Query() *Checker {
	length := len(c.providers)
	responses := make(chan int, length)
	c.wg.Add(length)
	go func() {
		c.positive = 0
		c.queried = 0
		for response := range responses {
			c.positive += response
			c.queried++
			c.wg.Done()
		}
	}()
	for _, provider := range c.providers {
		go func(provider string) {
			responses <- query(c, provider)
		}(provider)
	}
	c.length = length
	c.wg.Wait()
	close(responses)
	return c
}

// Stats returns the number of positive results along with the amount
// of blacklists supplied and the amount that were reachable.
// length is the number of providers
// queried is the number of providers who answered
// positive is the number of appearances reported
func (c *Checker) Stats() (positive, queried, length int) {
	return c.positive, c.queried, c.length
}

func (c *Checker) lookup(address string) ([]string, error) {
	return net.LookupHost(address)
}
