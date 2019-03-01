BUILD_VERSION   := $(shell cat version)

all:
	gox -osarch="darwin/amd64 linux/386 linux/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    	-ldflags "-w -s"

docker:
	docker build -t finab/bark-server:${BUILD_VERSION} .

clean:
	rm -rf dist

install:
	go install

.PHONY : all release docker clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://athens.azurefd.net
