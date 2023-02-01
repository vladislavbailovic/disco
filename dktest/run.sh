#!/bin/bash

set -e
go vet .

docker build . -t diskey:dktest

docker run --rm -d --name dktest-one diskey:dktest
docker run --rm -d --link dktest-one --name dktest-two diskey:dktest
docker run --rm -d --link dktest-one --name dktest-three diskey:dktest

sleep 15
echo "--- ONE ---"
docker logs dktest-one
echo "--- TWO ---"
docker logs dktest-two
echo "--- THREE ---"
docker logs dktest-three

docker stop dktest-one
docker stop dktest-two
docker stop dktest-three
