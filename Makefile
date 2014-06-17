.PHONY: test

test:
	for f in *.go; do \
		echo "Build $$f" ; \
		go build $$f ; \
	done

	echo "Build tavor with race detection"
	go install -race $(GOPATH)/src/github.com/zimmski/tavor/bin/tavor.go

	for f in *.tavor; do \
		echo "Testfuzz $$f" ; \
		tavor --input-file $$f ; \
	done
