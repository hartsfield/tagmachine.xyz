#!/bin/bash
pkill $1
go build -o $1
hmacss=$3 servicePort=$2 ./$1 & 
