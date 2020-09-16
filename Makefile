BUILD_VERSION   := $(shell cat version)
BUILD_DATE      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)

all:
	./cross_compile.sh

docker:
	cat deploy/Dockerfile | docker build -t finab/bark-server:${BUILD_VERSION} -f - .

release: clean all
	cp deploy/* dist
	ghr -u finb -t ${GITHUB_TOKEN} -replace -recreate --debug ${BUILD_VERSION} dist

clean:
	rm -rf dist

install:
	go install

.PHONY : all release docker clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.cn
