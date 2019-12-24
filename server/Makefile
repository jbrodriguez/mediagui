#
# Makefile to perform "live code reloading" after changes to .go files.
#
# To start live reloading run the following command:
# $ make serve
#

mb_version := $(shell cat ../VERSION)
mb_count := $(shell git rev-list HEAD --count)
mb_hash := $(shell git rev-parse --short HEAD)

# binary name to kill/restart
PROG = mediagui
 
# targets not associated with files
.PHONY: default build test coverage clean kill restart serve
 
# default targets to run when only running `make`
default: test
 
# clean up
clean:
	go clean

protobuf:
	protoc -I mediaagent/ mediaagent/agent.proto --go_out=plugins=grpc:mediaagent
 
# run formatting tool and build
build: clean
	go build fmt
	go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o ${PROG}

buildx: clean
	go build fmt
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o ${PROG}

agentx: clean
	env GOOS=linux GOARCH=amd64 go build -tags agent -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o agentx agent.go

agent: clean
	go build -tags agent -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o agentx agent.go

release: clean
	go build fmt
	go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o ${PROG}
	env GOOS=linux GOARCH=amd64 go build -tags agent -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o agentx agent.go
 
# run unit tests with code coverage
test: 
	go test -v
 
# generate code coverage report
coverage: test
	go build test -coverprofile=.coverage.out
	go build tool cover -html=.coverage.out
 
# attempt to kill running server
kill:
	-@killall -9 $(PROG) 2>/dev/null || true
 
# attempt to build and start server
restart:
	@make kill
	@make build; (if [ "$$?" -eq 0 ]; then (env GIN_MODE=debug ./${PROG} &); fi)

publish: build
	cp ./${PROG} ~/bin