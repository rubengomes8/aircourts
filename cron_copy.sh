#!/bin/bash

while true;
    echo starting courts finder
    do go run cmd/http/main.go -start=2022-05-31 -end=2022-05-31 -startTime=18:30 -indoor=false -maxStart=21:30 -slots=3 -indoor=false; 
    sleep 1800; # sleep 1 hour
done