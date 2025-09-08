# Bark

<img src="https://wx3.sinaimg.cn/mw690/0060lm7Tly1g0nfnjjxbbj30sg0sg757.jpg" width=200px height=200px />

[Bark](https://github.com/Finb/Bark) is an iOS App which allows you to push customed notifications to your iPhone.


## Table of Contents

   * [Bark](#bark)
      * [Installation](#installation)
         * [For Docker User](#for-docker-user)
         * [For General User](#for-general-user)
         * [For Developer](#for-developer)
         * [Nginx Proxy](#nginx-proxy)
      * [API V2](#api-v2)
      * [Other](#other)
         * [中文](#中文)
         * [Markdown Support](#markdown-support)         
      * [Contributing to bark-server](#contributing-to-bark-server)
         * [Development environment](#development-environment)
      * [Update](#update)


## Installation

### For Docker User

![Docker Automated build](https://img.shields.io/docker/automated/finab/bark-server.svg) ![Image Size](https://img.shields.io/docker/image-size/finab/bark-server?sort=date) ![License](https://img.shields.io/github/license/finb/bark-server)

The docker image is already available, you can use the following command to run the bark server:

``` sh
docker run -dt --name bark -p 8080:8080 -v `pwd`/bark-data:/data finab/bark-server
```

If you use the docker-compose tool, you can copy docker-copose.yaml under this project to any directory and run it:

``` sh
mkdir bark-server && cd bark-server
curl -sL https://github.com/Finb/bark-server/raw/master/deploy/docker-compose.yaml > docker-compose.yaml
docker compose up -d
```

### For General User 

- 1、Download precompiled binaries from the [releases](https://github.com/Finb/bark-server/releases) page
- 2、Add executable permissions to the bark-server binary: `chmod +x bark-server`
- 3、Start bark-server: `./bark-server --addr 0.0.0.0:8080 --data ./bark-data`
- 4、Test the server: `curl localhost:8080/ping`

**Note: Bark-server uses the `/data` directory to store data by default. Make sure that bark-server has permission to write to the `/data` directory, otherwise use the `-d` option to change the directory.**

### For Developer

Developers can compile this project by themselves, and the dependencies required for compilation:

- Golang 1.18+
- Go Mod Enabled(env `GO111MODULE=on`)
- Go Mod Proxy Enabled(env `GOPROXY=https://goproxy.cn`)
- [go-task](https://taskfile.dev/installation/) Installed

Run the following command to compile this project:

```sh
# Cross compile all platforms
task

# Compile the specified platform (please refer to Taskfile.yaml)
task linux_amd64
task linux_amd64_v3
```

**Note: The linux amd64 v3 architecture was added in go 1.18, see [https://github.com/golang/go/wiki/MinimumRequirements#amd64](https://github.com/golang/go/wiki/MinimumRequirements#amd64)**

### Use MySQL instead of Bbolt

Just run the server with `-dsn=user:pass@tcp(mysql_host)/bark`, it will use MySQL instead of file database Bbolt

## API V2

Please read [API_V2.md](docs/API_V2.md).

## Other

### 中文

- [https://bark.day.app/#/deploy](https://bark.day.app/#/deploy)

### Markdown support:

- [https://github.com/adams549659584/bark-server](https://github.com/adams549659584/bark-server)

## Contributing to bark-server

### Development environment

This project requires at least the golang 1.12 version to compile and requires Go mod support.

- Golang 1.16
- GoLand 2020.3 or other Go IDE
- Docker(Optional)

## Update 

Now the push certificate embedded in the program will never expire. You only need to update the program if the push fails due to the expired certificate.
