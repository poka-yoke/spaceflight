# capcom

## Description

`capcom` is a communication management tool.
Currently, it is mostly usable for AWS EC2 Security groups.

## Installation

    go get github.com/poka-yoke/spaceflight/mcc/capcom

## Usage

    capcom help
    cacpom list
    capcom add --source 198.234.12.34 sg-459d024
    capcom revoke --source 198.234.12.34 sg-459d024

## Name reasoning

It is called after the [CAPCOM](https://en.wikipedia.org/wiki/Flight_controller#Capsule_Communicator_.28CAPCOM.29) flight controller console.
