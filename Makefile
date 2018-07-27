PKG_NAME=byte

default: clean build

build:
	go install

clean:
	rm -f $(GOPATH)/bin/$(PKG_NAME)
