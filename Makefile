#!/usr/bin/make -f

test: fmt
	go test ./...

fmt:
	go mod tidy && go fmt ./...

docs:
	-go run *.go -help 2>&1 >/dev/null | grep -v 'exit status 2' > README.md

install:
	go install

.PHONE: test fmt docs install package