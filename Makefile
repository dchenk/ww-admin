SHELL := /bin/bash

.PHONY: fmt gen clean go_test front_test test

fmt:
	gofmt -s -w -e ./
	find . -type f -name "*.proto" | xargs clang-format -verbose -style file -i

# Generate Protocol Buffers code.
gen:
	for i in `find types -type f -name "*.proto"`; do \
		protoc --proto_path=types --go_out=plugins=grpc,paths=source_relative:types "$$i"; done

clean:
	rm types/*.pb.go

go_test:
	go test ./...

front_test:
	cd front && npm run test -- --coverage

test: go_test front_test
