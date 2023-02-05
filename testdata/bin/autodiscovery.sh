#!/bin/bash

set -e
go vet .

docker build . -t disco:autodiscovery -f testdata/Dockerfile-autodiscovery

docker run --rm -d --name autodiscovery-one disco:autodiscovery
docker run --rm -d --link autodiscovery-one --name autodiscovery-two disco:autodiscovery
docker run --rm -d --link autodiscovery-one --name autodiscovery-three disco:autodiscovery

sleep 20
docker stop autodiscovery-three

sleep 15
echo
echo "--- ONE ---"
docker logs autodiscovery-one
echo
echo "--- TWO ---"
docker logs autodiscovery-two


docker stop autodiscovery-one
docker stop autodiscovery-two

# sleep 15
# echo "--- ONE ---"
# docker logs autodiscovery-one
# echo "--- TWO ---"
# docker logs autodiscovery-two
# echo "--- THREE ---"
# docker logs autodiscovery-three


# docker stop autodiscovery-one
# docker stop autodiscovery-two
# docker stop autodiscovery-three
