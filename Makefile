build-all: build-re2 build-deduplicate-lines build-deduplicate-csv

build-re2:
	mkdir -p bin
	go build -o bin/re2 re2/*.go

build-deduplicate-lines:
	mkdir -p bin
	go build -o bin/deduplicate-lines deduplicate-lines/*.go

build-deduplicate-csv:
	mkdir -p bin
	go build -o bin/deduplicate-csv deduplicate-csv/*.go
