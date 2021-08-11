.PHONY: lint vendor clean

export GO111MODULE=on

default: lint 

lint:
	golangci-lint run

yaegi_test:
	yaegi test .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor
