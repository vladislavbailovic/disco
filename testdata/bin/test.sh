#!/bin/bash

set -e

go run . dispatch 3330 &
go run . storage 6661 &
go run . storage 6662 &
go run . storage 6663 &
go run . storage 6664 &

sleep 1

function check {
	echo
	echo -n "$1: "; curl -v http://localhost:3330/$1
	echo "---------------------------------"
	echo -n "$1: "; curl -X POST http://localhost:3330/$1 -d '{"test": "nanana"}'
	echo "---------------------------------"
	echo -n "$1: "; curl http://localhost:3330/$1
	echo
}

for x in {a..z}; do
	check $x
done

wait
