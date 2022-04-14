#!/bin/bash

while true;
    echo starting courts finder
    do ./bin/find_courts -email=true; 
    sleep 3600; # sleep 1 hour
done