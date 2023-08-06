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

### Nginx Proxy

Most users want to deploy the bark server on the intranet server, and then use Nginx to reverse proxy the bark server;
here is a simple Nginx configuration example (we assume that the bark server is listening at `192.168.1.123:8080`)

```sh
# generated 2020-03-26, Mozilla Guideline v5.4, nginx 1.17.7, OpenSSL 1.1.1d, modern configuration
# https://ssl-config.mozilla.org/#server=nginx&version=1.17.7&config=modern&openssl=1.1.1d&guideline=5.4
server {
    listen 80;
    listen [::]:80;
    # Replace bark.app.dev with your real domain name.
    server_name bark.app.dev;

    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    # Replace bark.app.dev with your real domain name.
    server_name bark.app.dev;

    ssl_certificate /path/to/signed_cert_plus_intermediates;
    ssl_certificate_key /path/to/private_key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:MozSSL:10m;  # about 40000 sessions
    ssl_session_tickets off;

    # modern configuration
    ssl_protocols TLSv1.3;
    ssl_prefer_server_ciphers off;

    # HSTS (ngx_http_headers_module is required) (63072000 seconds)
    add_header Strict-Transport-Security "max-age=63072000" always;

    # OCSP stapling
    ssl_stapling on;
    ssl_stapling_verify on;

    # verify chain of trust of OCSP response using Root CA and Intermediate certs
    ssl_trusted_certificate /path/to/root_CA_cert_plus_intermediates;

    # replace with the IP address of your resolver
    #resolver 127.0.0.1;

    location / {

        log_not_found on;
        # Replace http://192.168.1.123:8080 with the listening address of the bark server.
        proxy_pass http://192.168.1.123:8080;

        proxy_read_timeout 300;
        proxy_connect_timeout 300;
        proxy_redirect off;

        proxy_set_header Host              $host;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Real-IP         $remote_addr;

    }
}
```

### Use MySQL instead of Bbolt

Just run the server with `-dsn=user:pass@tcp(mysql_host)/bark`, it will use MySQL instead of file database Bbolt

## API V2

Please read [API_V2.md](docs/API_V2.md).

## Other

### 中文

- [https://day.app/2018/06/bark-server-document](https://day.app/2018/06/bark-server-document)

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
