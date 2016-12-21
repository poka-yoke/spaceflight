package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/poka-yoke/spaceflight/mcc/ttl"
)

type node struct {
	parent   *node
	children []*node
	content  *route53.ResourceRecordSet
}

func (n *node) isRoot() bool {
	if n.parent != nil {
		return false
	}
	return true
}

var zoneName string
var rootLookup map[*string]*node
var treeList []*node

// Init sets the flag parsing and input validations
func Init() {
	flag.StringVar(&zoneName, "zonename", "", "Hosted Zone's name to traverse")

	flag.Parse()

	if zoneName == "" {
		log.Fatal("Insufficient input parameters!")
	}
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("Failed to create session: %s", err)
		return
	}

	svc := route53.New(sess)
	zoneID := ttl.GetZoneID(zoneName, svc)

	records := ttl.GetResourceRecordSet(zoneID, svc)
	filteredRecords := ttl.FilterResourceRecordSetType(
		records,
		[]string{
			"A",
			"AAAA",
		},
	)
	// Roto. Necesita más filtrados.
	// 1. Identifica entradas relacionadas con el resto
	//     - Lookup table de val.Name -> val
	//     - Iterar para todas las entradas buscando aquellas cuyo
	//       val.ResourceRecords[0].Value aparezca como key en la LookupTable
	//     - Las relacionadas van a otra LookupTable
	// 2. Identifica las raíces entre las que están relacionadas
	//     - Se puede usar el algoritmo de más abajo
	// 3. Construye el árbol
	for _, val := range filteredRecords {
		root := &node{
			parent:  nil,
			content: val,
		}
		treeList = append(treeList, root)
		rootLookup[val.Name] = root
	}
	for _ = true; len(records) > 0; {
		for i, val := range records {
			if *val.Type == "CNAME" {
				root := rootLookup[val.ResourceRecords[0].Value]
				if root != nil {
					n := &node{
						content: val,
					}
					n.parent = root
					root.children = append(root.children, n)
					rootLookup[val.Name] = n
				}
			} else {
				records = append(records[:i], records[i+1:]...)
			}
		}
	}
}
