include .env

build:
	go build -v ${GOPACKAGES}

install:
	go install -v ${GOPACKAGES}

test:
	go test -v ${GOPACKAGES}

doc: doc/taocp-7-2-2-2-AlgorithmL.pdf

%.pdf: %.dot
	dot -Tpdf $< -o $@