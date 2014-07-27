.PHONY: test

test:
	for f in *.go; do \
		echo "Test binary $$f" ; \
		go build -o binout -race $$f ; \
		./binout ; \
		rm binout ; \
	done

	echo "Build tavor with race detection"
	go install -race $(GOPATH)/src/github.com/zimmski/tavor/bin/tavor.go

	for f in *.tavor; do \
		echo "Test format file $$f" ; \
		tavor --format-file $$f fuzz ; \
	done
