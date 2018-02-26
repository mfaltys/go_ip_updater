GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)


deps:
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/credentials
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/route53

run:
	go run \
		go_ip_updater/go_ip_updater.go \
		go_ip_updater/parselist.go \
    -loglevel debug

stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/go_ip_updater-$(GIT_HASH)-linux-amd64 go_ip_updater/*.go

clean:
	rm -rf bin/
