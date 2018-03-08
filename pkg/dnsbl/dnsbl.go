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

// Lookup contains the lookup function used
var Lookup = net.LookupHost
var wg sync.WaitGroup

// Stats of the DNSBL queries
var Stats struct {
	Length, Queried, Positive int
}

// Read introduces each line from io.Reader in a channel
func Read(in io.Reader) <-chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			wg.Add(1)
			out <- scanner.Text()
			Stats.Length++
		}
		close(out)
	}()
	return out
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

// Query queries a DNSBL and returns true if the argument gets a match
// in the BL.
func query(ipAddress, bl string, addresses chan<- int) {
	reversedIPAddress := fmt.Sprintf(
		"%v.%v",
		reverseAddress(ipAddress),
		bl,
	)
	result, _ := Lookup(reversedIPAddress)
	if len(result) > 0 {
		log.Printf("%v present in %v(%v)", reversedIPAddress, bl, result)
	}
	addresses <- len(result)
	Stats.Queried++
}

// Queries handles concurrency for Query. WaitGroup elements are added
// when reading the input
func Queries(ipAddress string, list <-chan string) {
	responses := make(chan int)
	for l := range list {
		go query(ipAddress, l, responses)
	}
	go func() {
		for response := range responses {
			if response > 0 {
				Stats.Positive += response
			}
			wg.Done()
		}
	}()
	wg.Wait()
	close(responses)
}
