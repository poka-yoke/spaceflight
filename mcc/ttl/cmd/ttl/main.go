package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/poka-yoke/spaceflight/mcc/ttl"
)

var zoneName string
var myttl int64
var verbose bool
var wait bool
var dryrun bool
var entryTypeFlag ttl.Filter
var entryNameFlag ttl.Filter

// Init sets the flag parsing and input validations
func Init() {
	flag.StringVar(&zoneName, "zonename", "", "Hosted Zone's name to traverse")
	flag.Int64Var(&myttl, "ttl", 300, "Desired TTL value")
	flag.BoolVar(&verbose, "v", false, "Increments output")
	flag.BoolVar(&wait, "w", false, "Waits for changes to complete")
	flag.BoolVar(&dryrun, "dryrun", false, "Does not commit the changes")
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
	list = ttl.FilterResourceRecordSetType(list, entryTypeFlag)
	list, list2 := ttl.SplitResourceRecordSetTypeOnNames(list, entryNameFlag)
	if len(list2) > 0 {
		list = list2
	}
	changeResponse, err := ttl.UpsertResourceRecordSetTTL(list, myttl, &zoneID, svc)
	if err != nil {
		log.Panic(err.Error())
	}
	if wait && !dryrun {
		ttl.WaitForChangeToComplete(changeResponse.ChangeInfo, svc)
	}
}
