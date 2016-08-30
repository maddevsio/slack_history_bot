TARGET=history_bot

all: fmt clean build

clean:
	rm -rf $(TARGET)

depends:
	go get -u -v -tags 'goleveldb libstemmer' ./...

build:
	go build -tags 'goleveldb libstemmer'  -v -o  $(TARGET) main.go

fmt:
	go fmt ./...

test:
	go test -v ./...
