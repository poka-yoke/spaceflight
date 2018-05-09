# trek

## Description

`trek` is a redirection generation tool.
Currently, it only works for Nginx.

## Installation

    go get github.com/poka-yoke/spaceflight/mcc/trek

## Usage

    trek help
    cat nginx.conf | trek add --original help.example.com --final https://www.example.com/help.html

## Name reasoning

It is called after the [Trek](https://www.nasa.gov/sites/default/files/atoms/files/g-28367c_trek.pdf) software.
