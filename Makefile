PKG_NAME=byte
TEST_PKGS=./utils

default: test build

build:
	go install

pretest: clean
	mkdir coverage

test: pretest
	go test $(TEST_PKGS) -v -coverprofile=coverage/coverage.out

clean:
	rm -rf $(GOPATH)/bin/$(PKG_NAME) ./coverage

reset:
	rm ./test/dest/*
