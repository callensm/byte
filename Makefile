PKG_NAME=byte
TEST_PKGS=./utils

default: test build

build:
	go install

test: clean
	go test $(TEST_PKGS) -v -cover

clean:
	rm -f $(GOPATH)/bin/$(PKG_NAME)

reset:
	rm ./test/dest/*
