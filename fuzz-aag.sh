#!/bin/sh

mkdir out

for V in {1..1000}
do
	echo "Run #$V"

	tavor --format-file aag.tavor --verbose fuzz --exec "bin/aigtoaig" --exec-argument-type "stdin" --exec-exact-exit-code 0 --result-folder out/ --exit-on-error --exec-do-not-remove-tmp-files-on-error
done

