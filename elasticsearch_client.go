package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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
	log.Println("Starting es-manager")

	var message interface{}
	if len(os.Args) == 2 {
		err := json.Unmarshal([]byte(os.Args[1]), &message)
		if err != nil {
			log.Panicln(err)
		}
		log.Println(message)
		// TODO: Parse the message contents here
		for k, v := range message.(map[string]interface{}) {
			log.Println(k, v)
		}
	} else {
		log.Fatalln("Not enough Parameters!")
	}

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	h := http.DefaultClient

	server := "http://devex-es-develop-classic.us-east-1.elasticbeanstalk.com:80"
	log.Printf("Connecting to %s", server)
	client, err := elastic.NewClient(h, server)
	if err != nil {
		// Handle error
		log.Panicln(err)
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://devex-es-develop-classic.us-east-1.elasticbeanstalk.com:80")
	if err != nil {
		// Handle error
		log.Panicln(err)
	}
	log.Printf("Elasticsearch version %s\n", esversion)

	ret, err := updateSettings(client, h, "/_all/_settings", "auto_expand_replicas: 1-all")
	if err != nil {
		// Handle error
		log.Panicln(err)
	}
	if ret != 200 {
		log.Printf("Error! %d", ret)
	} else {
		log.Println("All good!")
	}
	ret, err = updateSettings(client, h, "/_cluster/settings", "{ \"transient\": { \"discovery.zen.minimum_master_nodes\": 2 } }")
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
