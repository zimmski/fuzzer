#!/bin/sh

for f in *.go
do
	echo "Test binary $f"

	go build -o binout -race $f
	if [ $? -ne 0 ]; then
		exit
	fi

	./binout
	if [ $? -ne 0 ]; then
		exit
	fi

	rm binout
done

echo "Build tavor with race detection"
make -C $GOPATH/src/github.com/zimmski/tavor debug-install

for f in *.tavor
do
	echo "Test format file $f"

	tavor --format-file $f fuzz
	if [ $? -ne 0 ]; then
		exit
	fi
done
