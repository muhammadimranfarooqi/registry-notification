VERSION=`git rev-parse --short HEAD`
flags=-ldflags="-s -w -X main.version=${VERSION}"

all: vet build

vet:
	go vet .

build:
	go clean; rm -rf pkg; go build -o registry-notification ${flags}

build_debug:
	go clean; rm -rf pkg; go build -o registry-notification ${flags} -gcflags="-m -m"

build_all: build_osx build_linux build

build_osx:
	go clean; rm -rf pkg registry-notification_osx; GOOS=darwin go build -o registry-notification ${flags}

build_linux:
	go clean; rm -rf pkg registry-notification_linux; GOOS=linux go build -o registry-notification ${flags}

build_power8:
	go clean; rm -rf pkg registry-notification_power8; GOARCH=ppc64le GOOS=linux go build -o registry-notification ${flags}

build_arm64:
	go clean; rm -rf pkg registry-notification_arm64; GOARCH=arm64 GOOS=linux go build -o registry-notification ${flags}

build_windows:
	go clean; rm -rf pkg registry-notification.exe; GOARCH=amd64 GOOS=windows go build -o registry-notification ${flags}

install:
	go install

clean:
	go clean; rm -rf pkg

test : test1

test1:
	go test -v -bench=.

