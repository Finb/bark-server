BUILD_VERSION   := $(shell cat version)
BUILD_DATE      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)

all:
	bash .cross_compile.sh bark-server

docker:
	cat deploy/Dockerfile | docker build -t finab/bark-server:${BUILD_VERSION} -f - .

buildx:
	bash .buildx.sh

release: clean all
	cp deploy/* dist
	ghr -u finb -t ${GITHUB_TOKEN} -replace -recreate -name "Bump ${BUILD_VERSION}" --debug ${BUILD_VERSION} dist

pre-release: clean all
	cp deploy/* dist
	ghr -u finb -t ${GITHUB_TOKEN} -replace -recreate -prerelease -name "Bump ${BUILD_VERSION}" --debug ${BUILD_VERSION} dist

clean:
	rm -rf dist

install:
	go install

.PHONY : all release docker buildx clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.cn
