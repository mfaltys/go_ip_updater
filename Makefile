GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)


deps:
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/credentials
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/route53
	go get github.com/unixvoid/glogger

run:
	go run \
		go_ip_updater/go_ip_updater.go \
		go_ip_updater/parselist.go \
    -loglevel debug

prep_aci: stat
	mkdir -p go_ip_updater-layout/rootfs/deps/
	cp deps/manifest.json go_ip_updater-layout/manifest
	cp bin/go_ip_updater* go_ip_updater-layout/rootfs/go_ip_updater

build_aci: prep_aci
	actool build go_ip_updater-layout go_ip_updater.aci
	@echo "go_ip_updater.aci built"

build_travis_aci: prep_aci
	wget https://github.com/appc/spec/releases/download/v0.8.7/appc-v0.8.7.tar.gz
	tar -zxf appc-v0.8.7.tar.gz
	# build image
	appc-v0.8.7/actool build go_ip_updater-layout go_ip_updater.aci && \
	rm -rf appc-v0.8.7*
	@echo "go_ip_updater.aci built"

stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/go_ip_updater-$(GIT_HASH)-linux-amd64 go_ip_updater/*.go

clean:
	rm -rf bin/
	rm -f go_ip_updater.aci
	rm -rf go_ip_updater-layout
