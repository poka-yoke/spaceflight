package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/Devex/spaceflight/pkg/roosa"
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
	zoneID := roosa.GetZoneID(zoneName, svc)

	referenceTreeList := roosa.NewReferenceTreeList(
		roosa.GetResourceRecordSet(zoneID, svc),
	)
	fmt.Print(referenceTreeList)
}
