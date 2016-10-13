package main

import (
	//	"encoding/json"
	"fmt"
	"net/http"
	//	"time"

	"gopkg.in/olivere/elastic.v1"
)

func updateSettings(client *elastic.Client, h *http.Client, url string, setting string) (retcode int, err error) {
	req, err := client.NewRequest("PUT", url)
	if err != nil {
		return
	}
	req.SetBodyString(setting)
	res, err := h.Do((*http.Request)(req))
	if err != nil {
		return
	}
	defer res.Body.Close()
	retcode = res.StatusCode
	return
}

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

	ret, err := updateSettings(client, h, "/_all/_settings", "auto_expand_replicas: 1-all")
	if err != nil {
		// Handle error
		panic(err)
	}
	if ret != 200 {
		fmt.Printf("Error! %d", ret)
	} else {
		fmt.Println("All good!")
	}
}
