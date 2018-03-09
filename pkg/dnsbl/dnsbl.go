package dnsbl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// Checker controls the flow of package and provides a single point of
// entry for its users.
type Checker struct {
	length, queried, positive int

	// lookup contains the lookup function used
	lookup func(string) ([]string, error)
	// Functional control
	wg sync.WaitGroup
}

// NewChecker creates a new, default configured Checker
func NewChecker() *Checker {
	return &Checker{lookup: net.LookupHost}
}

// Query handles concurrency for Query. WaitGroup elements are added
// when reading the input
func (c *Checker) Query(ipAddress string, lists io.Reader) *Checker {
	list := c.read(lists)
	responses := make(chan int)
	for l := range list {
		go c.query(ipAddress, l, responses)
	}
	go func() {
		for response := range responses {
			if response > 0 {
				c.positive += response
			}
			c.wg.Done()
		}
	}()
	c.wg.Wait()
	close(responses)
	return c
}

// Stats returns the number of positive results along with the amount
// of blacklists supplied and the amount that were reachable.
func (c *Checker) Stats() (int, int, int) {
	return c.positive, c.queried, c.length
}

// read introduces each line from io.Reader in a channel
func (c *Checker) read(in io.Reader) <-chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			c.wg.Add(1)
			out <- scanner.Text()
			c.length++
		}
		close(out)
	}()
	return out
}

// Query queries a DNSBL and returns true if the argument gets a match
// in the BL.
func (c *Checker) query(ipAddress, bl string, addresses chan<- int) {
	reversedIPAddress := fmt.Sprintf(
		"%v.%v",
		reverseAddress(ipAddress),
		bl,
	)
	result, _ := c.lookup(reversedIPAddress)
	if len(result) > 0 {
		log.Printf("%v present in %v(%v)", reversedIPAddress, bl, result)
	}
	addresses <- len(result)
	c.queried++
}

// Reverse reverses slice of string elements.
func reverse(original []string) {
	for i := len(original)/2 - 1; i >= 0; i-- {
		opp := len(original) - 1 - i
		original[i], original[opp] = original[opp], original[i]
	}
}

// ReverseAddress converts IP address in string to reversed address for query.
func reverseAddress(ipAddress string) (reversedIPAddress string) {
	ipAddressValues := strings.Split(ipAddress, ".")
	reverse(ipAddressValues)
	reversedIPAddress = strings.Join(ipAddressValues, ".")
	return
}
