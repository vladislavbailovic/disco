#!/bin/bash

set -ex
go vet .

docker build . -t disco:storage -f testdata/Dockerfile-storage

docker run --rm -d --name storage-one disco:storage
docker run --rm -d --link storage-one --name storage-two disco:storage
docker run --rm -d --link storage-one --name storage-three disco:storage

sleep 15
docker stop storage-three

sleep 30
echo
echo "--- ONE ---"
docker logs storage-one
echo
echo "--- TWO ---"
docker logs storage-two


docker stop storage-one
docker stop storage-two

# sleep 15
# echo "--- ONE ---"
# docker logs storage-one
# echo "--- TWO ---"
# docker logs storage-two
# echo "--- THREE ---"
# docker logs storage-three


# docker stop storage-one
# docker stop storage-two
# docker stop storage-three
