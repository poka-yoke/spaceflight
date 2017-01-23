package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/Devex/spaceflight/mcc/ttl"
)

var zoneName string
var myttl int64
var verbose bool
var wait bool
var dryrun bool
var exclude bool
var entryTypeFlag ttl.Filter
var entryNameFlag ttl.Filter

// Init sets the flag parsing and input validations
func Init() {
	flag.StringVar(&zoneName, "zonename", "", "Hosted Zone's name to traverse")
	flag.Int64Var(&myttl, "ttl", 300, "Desired TTL value")
	flag.BoolVar(&verbose, "v", false, "Increments output")
	flag.BoolVar(&wait, "w", false, "Waits for changes to complete")
	flag.BoolVar(&dryrun, "dryrun", false, "Does not commit the changes")
	flag.BoolVar(
		&exclude,
		"exclude",
		false,
		"Filter lists are used to exclude records instead of including them",
	)
	flag.Var(&entryTypeFlag,
		"t",
		"Comma-separated list of entry record types")
	flag.Var(&entryNameFlag,
		"n",
		"Comma-separated list of entry record names")

	flag.Parse()

	if zoneName == "" {
		log.Fatal("Insufficient input parameters!")
	}

	if len(entryTypeFlag) == 0 {
		for _, f := range []string{
			"A",
			"AAAA",
			"CNAME",
			"MX",
			"NAPTR",
			"PTR",
			"SPF",
			"SRV",
			"TXT",
		} {
			entryTypeFlag = append(entryTypeFlag, f)
		}
	}
	ttl.Verbose = verbose
	ttl.Dryrun = dryrun
}

func filter(l []*route53.ResourceRecordSet) []*route53.ResourceRecordSet {
	list := l
	if exclude {
		log.Println("Exclude!")
		list = ttl.FilterResourceRecords(
			l,
			entryTypeFlag,
			func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
				if *elem.Name != filter {
					return elem
				}
				return nil
			},
		)
		if len(entryNameFlag) > 0 {
			list = ttl.FilterResourceRecords(
				list,
				entryNameFlag,
				func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
					if *elem.Name != filter {
						return elem
					}
					return nil
				},
			)
		}
	} else {
		list = ttl.FilterResourceRecords(
			l,
			entryTypeFlag,
			func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
				if *elem.Name == filter {
					return elem
				}
				return nil
			},
		)
		if len(entryNameFlag) > 0 {
			list = ttl.FilterResourceRecords(
				list,
				entryNameFlag,
				func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
					if *elem.Name == filter {
						return elem
					}
					return nil
				},
			)
		}
	}
	return list
}

func main() {
	Init()
	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("Failed to create session: %s", err)
		return
	}

	svc := route53.New(sess)
	zoneID := ttl.GetZoneID(zoneName, svc)

	list := ttl.GetResourceRecordSet(zoneID, svc)
	// Filter list in between
	list = filter(list)
	changeResponse, err := ttl.UpsertResourceRecordSetTTL(list, myttl, &zoneID, svc)
	if err != nil {
		log.Panic(err.Error())
	}
	if wait && !dryrun {
		ttl.WaitForChangeToComplete(changeResponse.ChangeInfo, svc)
	}
}
