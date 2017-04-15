# got

## Description

`got` is a generic DNS service management tool.
Basically, it allows you to create, modify, and delete DNS records.

It's not focused on management of DNS servers, only the service data.

## Installation

    go get github.com/poka-yoke/spaceflight/mcc/got

## Usage

    got help
    help ttl --zone example.com -ttl 30
    got upsert --name www.example.com. --zone example.com --ttl 300 --type CNAME myserver.example.com
    got ttl --zone example.com -ttl 360

## Name reasoning

It is called after [Seymour Liebergot](https://en.wikipedia.org/wiki/Seymour_Liebergot) who manned the [EECOM](https://en.wikipedia.org/wiki/Flight_controller#Electrical.2C_Environmental_and_Consumables_Manager_.28EECOM.29) flight controller console during Apolo XIII explosion, and who helped guiding the spaceship back to Earth.
