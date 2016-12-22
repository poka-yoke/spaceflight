package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/Devex/spaceflight/mcc/roosa"
	"github.com/poka-yoke/spaceflight/mcc/ttl"
)

var zoneName string

// Init sets the flag parsing and input validations
func Init() {
	flag.StringVar(&zoneName, "zonename", "", "Hosted Zone's name to traverse")

	flag.Parse()

	if zoneName == "" {
		log.Fatal("Insufficient input parameters!")
	}
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

	records := ttl.GetResourceRecordSet(zoneID, svc)
	referenceTrees := roosa.GetReferenceTrees(records)
	roosa.PrintReferenceTrees(referenceTrees)
}
