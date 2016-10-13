#!/bin/bash

GOOS=linux GOARCH=amd64 go build elasticsearch_client.go
zip elasticsearch_client elasticsearch_client main.js
