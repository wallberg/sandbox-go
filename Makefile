include .env

build:
	go build -v ${GOPACKAGES}

install:
	go install -v ${GOPACKAGES}

test:
	go test -v ${GOPACKAGES}
