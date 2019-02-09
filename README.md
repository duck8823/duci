# duci
![Language](https://img.shields.io/badge/language-go-74CCDC.svg) ![GitHub release](https://img.shields.io/github/release/duck8823/duci.svg?colorB=7E7E7E) [![GoDoc](https://godoc.org/github.com/duck8823/duci?status.svg)](https://godoc.org/github.com/duck8823/duci) [![Build Status](https://travis-ci.org/duck8823/duci.svg?branch=master)](https://travis-ci.org/duck8823/duci) [![Coverage Status](https://coveralls.io/repos/github/duck8823/duci/badge.svg?branch=master)](https://coveralls.io/github/duck8823/duci?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/duck8823/duci)](https://goreportcard.com/report/github.com/duck8823/duci) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

duci \[zushi\] (Docker Under Continuous Integration) is a simple ci server.  
 

## DSL is Unnecessary For CI
Let's define the task in the task runner.  
Let's define the necessary infrastructure for the task in the Dockerfile.  
duci just only execute the task in docker container.  

## Features
- Execute the task in Docker container
- Execute the task triggered by GitHub pull request comment or push 
- Execute tasks asynchronously
- Create GitHub commit status
- Store and Show logs

## How to use
### Target Repository
The target repository must have Dockerfile in repository root or `.duci/Dockerfile`.  
If there is `.duci/Dockerfile`, duci read it preferentially.  
  
In Dockerfile, I suggest to use `ENTRYPOINT`.

e.g.
```Dockerfile
ENTRYPOINT ["mvn"]
CMD ["compile"]
```

```Dockerfile
ENTRYPOINT ["fastlane"]
CMD ["build"]
```

When push to github, duci execute `mvn compile` / `fastlane build`.  
And when comment `ci test` on github pull request, execute `mvn test` / `fastlane test`.  

### Using host environment value
If exists `ARG` instruction in Dockerfile, override environment variable from host.  

```Dockerfile
ARG FOO=default
ARG BAR
```

### Runtime configuration
#### volumes
You can use volumes options for external dependency, cache and etc.  
Set configurations in `.duci/config.yml`  

```yaml
volumes:
  - '/path/to/host/dir:/path/to/container/dir'
```

#### environment variable
You can set environment variables in docker container.  
Add the following to `.duci/config.yml`

```yaml
environments:
  - ENVIRONMENT_VAIRABLE=value
```

## Server Settings
### Installation
Download the asset for your platform from [releases page](https://github.com/duck8823/duci/releases).  
Extract the downloaded file, and add the binary to your `PATH`. 

### Setting SSH (optional)
If target repository is private, You can use SSH key to clone repository from github.com.  
Please set the public key of the pair at https://github.com/settings/keys.

### Add Webhooks to Your GitHub repository
duci start to listen webhook with port `8080` (default) and endpoint `/`.  
In GitHub target repository settings (`https://github.com/<owner>/<repository>/settings/hooks`),
Add endpoint of duci to `Payload URL` and `application/json` to `Content type` respectively.

### Run Server
```bash
$ duci server
```

### Server Configuration file
You can specify configuration file with `-c` option.
The configuration file must be yaml format.
Possible values ​​are as follows.

```yaml
server:
  workdir: '/path/to/tmp/duci'
  port: 8080
  database_path: '$HOME/.duci/db'
github:
  # (optional) You can use SSH key to clone. ex. '${HOME}/.ssh/id_rsa'
  ssh_key_path: ''
  # For create commit status. You can also use environment variable
  api_token: ${GITHUB_API_TOKEN}
job:
  timeout: 600
  concurrency: 4 # default is number of cpu
```

You can check the configuration values.

```bash
$ duci config
```

## Using Docker
You can use Docker to run server.
```
$ docker run -p 8080:8080 \
             -e GITHUB_API_TOKEN=<your toekn> \
             -v /var/run/docker.sock:/var/run/docker.sock \
             duck8823/duci
```

When you want to clone with SSH in container,
```
$ docker run -p 8080:8080 \
             -e GITHUB_API_TOKEN=<your toekn> \
             -e SSH_KEY_PATH=/root/.ssh/id_rsa \
             -v ~/.ssh:/root/.ssh:ro \ 
             -v /var/run/docker.sock:/var/run/docker.sock \
             duck8823/duci
```

## Run with docker-compose
With docker-compose, you can also start ui and reverse proxy together.

```bash
$ git clone https://github.com/duck8823/duci.git
$ cd duci
$ docker-compose up -d
```

If you start up on another host, set your host name (default: `localhost`) to environment variable `DUCI_HOST`.

## Read job log
GitHub send payload as webhook including `X-GitHub-Delivery` header.  
You can read job log with the `X-GitHub-Delivery` value formatted UUID.

```bash
$ curl -XGET http://localhost:8080/logs/{X-GitHub-Delivery}
```

The endpoint returns NDJSON (Newline Delimited JSON) formatted log.

```jsons
{"time":"2018-09-21T22:19:42.572879+09:00","message":"Step 1/10 : FROM golang:1.11-alpine"}
{"time":"2018-09-21T22:19:42.573093+09:00","message":"\n"}
{"time":"2018-09-21T22:19:42.573494+09:00","message":" ---\u003e 233ed4ed14bf\n"}
{"time":"2018-09-21T22:19:42.573616+09:00","message":"Step 2/10 : MAINTAINER shunsuke maeda \u003cduck8823@gmail.com\u003e"}
{"time":"2018-09-21T22:19:42.573734+09:00","message":"\n"}
...
```

## Health Check
This server has an health check API endpoint (`/health`) that returns the health of the service. The endpoint returns `200` status code if all green.  

```bash
$ curl -XGET -I http://localhost:8080/health
```

```
HTTP/1.1 200 OK
Date: Wed, 31 Oct 2018 20:33:42 GMT
Content-Length: 0
```

The check items are as follows

- Whether the Docker daemon is running or not

You can also check with `health` sub-command.

```bash
$ duci health
```

```
INFO[14/Jan/2019 07:17:38.864] ok.
```

## License
MIT License

Copyright (c) 2018 Shunsuke Maeda

See [LICENSE](./LICENSE) file
