#-include .env

PROJECT_NAME := $(shell basename "$(PWD)")

PROJ_BASE := $(shell pwd -LP)
PROJ_BUILD_PATH := $(PROJ_BASE)/build

# PID file will keep the process id of the server
PID := /tmp/.$(PROJECTNAME).pid

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECT_NAME)-stderr.txt
# Redirect error output to a file, so we can show it in development mode.
STDOUT := /tmp/.$(PROJECT_NAME)-stdout.txt

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

## init: Simple initialization. Make `third_party/protoc-gen.sh` executable
init:
	@echo " > Simple initialization"
	@echo " >> Make `third_party/protoc-gen.sh` executable"
	@chmod +x $(PROJ_BASE)/third_party/protoc-gen.sh

## clean: Clean build files.
#clean:
#	@echo "  >  Clean build files. Runs `go clean` internally."
#	@-$(MAKE) clean-build clean-proto go-clean
clean:
	@echo "  >  Clean build files."
	@-$(MAKE) clean-build clean-proto

clean-build:
	@echo "  >  Clean build"
	@-rm $(PROJ_BUILD_PATH)/* 2> /dev/null

clean-proto:
	@echo "  >  Clean proto"
	@-rm -f $(PROJ_BASE)/api/swagger/**/*.json 2> /dev/null
	@-rm -f $(PROJ_BASE)/pkg/api/v1/*.pb.* 2> /dev/null
	@echo

## gen-proto: Generate proto
gen-proto:
	@echo "  >  Generate proto"
	@$(shell $(PROJ_BASE)/third_party/protoc-gen.sh)

## build-all: Runs `gen-proto` `build-server` `c`
build-all:
	@echo "  >  Build all"
	@-$(MAKE) gen-proto build-server

## run-server: RUN_OPTIONS='-grpc-port= -db-host= --db-port= -db-user= -db-password= -db-name='
run-server:
	@echo "  >  Running server"
	@$(PROJ_BUILD_PATH)/server $(RUN_OPTIONS)

start-server: stop-server
	@echo "  >  starting grpc server"
	@-$(PROJ_BUILD_PATH)/server $(RUN_OPTIONS) 2>$(STDOUT) & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"
	@echo "  >  stoud at $(STDOUT)"

stop-server:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)


logs:
	@tail -f -n 100 $(STDOUT)

## build-server: Build grpc-server
build-server:
	@echo "  >  Build server"
	@go build -o $(PROJ_BUILD_PATH)/server $(PROJ_BASE)/cmd/server/main.go

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@go get $(get)

#go-clean:
#	@echo "  >  Cleaning build cache"
#	@go clean $(PROJ_BUILD_PATH)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo