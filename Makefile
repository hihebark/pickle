TARGET=pickle

all: deps build

deps: godep
	@dep ensure

build:
	@go build -o $(TARGET) .

clean:
	@rm -rf $(TARGET)
	@rm -rf build

install:
	@cp $(TARGET) /usr/local/bin/

godep:
	@go get -u github.com/golang/dep/...
