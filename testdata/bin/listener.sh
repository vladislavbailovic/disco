#!/bin/bash

function finish() {
	echo "Well this is now done, kthxbai"
	docker stop listener-three
	docker stop listener-two
	docker stop listener-one
}
trap finish SIGINT

go vet .

docker build . -t disco:listener -f testdata/Dockerfile-listener

docker run --rm -d --name listener-one -p3366:6660 disco:listener
docker run --rm -d --link listener-one --name listener-two disco:listener
docker run --rm -d --link listener-one --name listener-three disco:listener

docker logs --follow listener-one&
docker logs --follow listener-two&
docker logs --follow listener-three&

wait
