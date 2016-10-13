package main

import (
	//	"encoding/json"
	"fmt"
	"net/http"
	//	"time"

	"gopkg.in/olivere/elastic.v1"
)

func main() {
	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	h := http.DefaultClient
	client, err := elastic.NewClient(h, "http://devex-es-develop-classic.us-east-1.elasticbeanstalk.com:80")
	if err != nil {
		// Handle error
		panic(err)
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://devex-es-develop-classic.us-east-1.elasticbeanstalk.com:80")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("_settings").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if exists {
		fmt.Println("It exists")
	} else {
		fmt.Println("It doesn't exist")
	}
	url := "/_all/_settings"
	req, err := client.NewRequest("PUT", url)
	if err != nil {
		panic(err)
	}
	req.SetBodyString("auto_expand_replicas: 1-all")
	_, err = h.Do((*http.Request)(req))
	if err != nil {
		panic(err)
	}

}
