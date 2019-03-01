# Bark

<img src="https://wx3.sinaimg.cn/mw690/0060lm7Tly1g0nfnjjxbbj30sg0sg757.jpg" width=200px height=200px />

[Bark](https://github.com/Finb/Bark) is an iOS App which allows you to push customed notifications to your iPhone.

## Installation

### For Docker User

![Docker Automated build](https://img.shields.io/docker/automated/finab/bark-server.svg) ![MicroBadger Size](https://img.shields.io/microbadger/image-size/finab/bark-server.svg) ![MicroBadger Layers](https://img.shields.io/microbadger/layers/finab/bark-server.svg)

The docker image is already available, you can use the following command to run the bark server:

``` sh
docker run -dt --name bark -p 8080:8080 -v `pwd`/bark-data:/data finab/bark-server
```

If you use the docker-compose tool, you can copy docker-copose.yaml under this project to any directory and run it:

``` sh
mkdir bark && cd bark
curl -sL https://git.io/fhAsj > docker-compose.yaml
docker-compose up -d
```

### For General User 

- 1、Download precompiled binaries from the release page
- 2、Add executable permissions to the bark-server binary: `chmod +x bark-server`
- 3、Start bark-server: `./bark-server -l 0.0.0.0 -p 8080 -d ./bark-data`
- 4、Test the server: `curl localhost:8080/ping`

**Note: Bark-server uses the /data directory to store data by default. Make sure that bark-server has permission to write to the /data directory, otherwise use the `-d` option to change directories.**

#### Other documents:

- [https://day.app/2018/06/bark-server-document/](https://day.app/2018/06/bark-server-document/)
  
## Contributing to bark-server

### Development environment

This project requires at least the golang 1.12 version to compile and requires Go mod support.

- Golang 1.12
- GoLand 2018.3.4 or other Go IDE
- Docker(Optional)

## Update 

The push certificate embedded in the program expires on **`2020/01/30`**, please update the program after **`2019/12/01`**
