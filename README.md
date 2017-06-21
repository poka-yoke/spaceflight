# Spaceflight project [![Build Status](https://travis-ci.org/poka-yoke/spaceflight.svg?branch=master)](https://travis-ci.org/poka-yoke/spaceflight) [![Code Climate](https://codeclimate.com/github/poka-yoke/spaceflight/badges/gpa.svg)](https://codeclimate.com/github/poka-yoke/spaceflight) [![Coverage Status](https://coveralls.io/repos/github/poka-yoke/spaceflight/badge.svg?branch=master)](https://coveralls.io/github/poka-yoke/spaceflight?branch=master)

## Aim

The Spaceflight project aims to build a toolbox to ease administration and orchestration of infrastructure, in the broadest sense of those concepts, with the goal of managing human errors.                                                                        

## Definition

This project is focused to create a set of composable tools and components, independent and scoped, mostly built with Go, and intended to be useful for managing AWS, but looking forward to see other platforms and providers contributed.                         

You'll find we talk about tool types. Keep reading for a description:

* Mission Control Consoles (MCC): These are tools thought to be executed from administration nodes, i.e. your computer.
* Mobile Servicing Systems (MSS): These are tools thought to be executed under certain automated systems, like checks in a monitoring system.
* Orbiters: These are tools thought to be executed in serverless infrastructure.

Finally, here you can find a relation of tools:

| Tool | Type | Description |
|---|---|---|
| [got](mcc/got) | MCC | Tool to manage changes in DNS |
| [capcom](mcc/capcom) | MCC | Communicactions management tool |
| [roosa](mcc/roosa) | MCC | Detection and visualization of relationships |
| [dextre](mss/dextre) | MSS | Spam blacklists presence and mail reputation check |
| [trek](mcc/trek) | MCC | Redirection management tool |
| [fido](mcc/fido) | MCC | Configuration management tool |
| [surgeon](mcc/surgeon) | MCC | (WIP) Service Diagnosing tool |
| [es](orbiter/es) | Orbiter | (WIP) ElasticSearch automatic tuning Lambda |

## Collaborations

In case you find a bug, or you need something to be added, you can post an Issue, and your need will be discussed.
If you want to collaborate with simple fixes or additions through personal forks, you can do it by issuing a PR.
For collaborating through corporate or team forks, let us know your interest to find a suitable one.

### PR acceptation policy

In order to have a PR accepted into spaceflight, it must:

* Include as many tests as needed to keep the code maintainable, and pass previous tests.
* All functions and types added must be single concern.
* Composability must be kept.
* Go fmt is a must. Besides that, we're looking into tools like go-simple and staticcheck.

We'll not be accepting PRs if they:

* Break use cases.
* Don't follow existing philosophy.

## Licensing and philosophy

Our philosophy relies on free software and open source, and we want the code in this repository to be shared with everyone, receiving contributions from anyone. So we agreed on these licensing:                                                                                                 

- All the tools (i.e. commands, Go’s main package) in this repository are licensed under GPL License, so these can’t be imported in to closed projects.                                                                                                              
- All the libraries are licensed under 3-Clause BSD License, to foster sharing.

Of course, we hope to hear back from you, whatever you do with these tools.
