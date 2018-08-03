PKG_NAME=byte
TEST_PKGS=./utils
TEST_OUT_DIR=coverage
COVER_FILE=output

default: test build

build:
	go install

test: clean pretest
	go test $(TEST_PKGS) -coverprofile $(COVER_FILE) -outputdir $(TEST_OUT_DIR)

pretest:
	mkdir coverage

travis: clean
	go test $(TEST_PKGS) -v -cover

clean:
	rm -rf $(GOPATH)/bin/$(PKG_NAME) ./coverage

reset:
	rm ./test/dest/*
