#!/bin/bash

if [[ $# -ne 1 ]]; then
  echo "Need path to directory or rjson file with raw sample data"
  exit 1
fi

echo "In browser open http://live.wearkinetic.com"

filepath=$1

go run gothrough.go -fp $filepath -speed 1 | go run /Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/stater/cmd/stater/*go -lift -lift.send.point | go run /Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/transporter/*go -ct=ipc -st=ws -a=point
