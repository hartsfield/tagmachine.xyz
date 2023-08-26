#!/bin/bash
pkill -9 $1
go build -o $1
hmacss=$3 servicePort=$2 logFilePath=./logfile.txt ./$1 & 
