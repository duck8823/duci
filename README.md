# duci
![Language](https://img.shields.io/badge/language-go-74CCDC.svg)
![GitHub release](https://img.shields.io/github/release/duck8823/duci.svg?colorB=7E7E7E)
[![Build Status](https://travis-ci.org/duck8823/duci.svg?branch=master)](https://travis-ci.org/duck8823/duci)
[![Coverage Status](https://coveralls.io/repos/github/duck8823/duci/badge.svg?branch=master)](https://coveralls.io/github/duck8823/duci?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/duck8823/duci)](https://goreportcard.com/report/github.com/duck8823/duci)
[![codebeat badge](https://codebeat.co/badges/dfae99c0-e051-4baa-b693-7869cc25069b)](https://codebeat.co/projects/github-com-duck8823-duci-master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

duci \[zushi\] (<u>D</u>ocker <u>U</u>nder <u>C</u>ontinuous <u>I</u>ntegration) is a small ci server.  
 

## DSL is Unnecessary For CI
Let's define the task in the task runner.  
In the Dockerfile, let's define the necessary infrastructure for the task.  
duci just execute the task in docker container. 

## Features
- Run task in Docker container
- The task is triggered by pull request comment and push 
- Create GitHub commit status
- Execute tasks asynchronously

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

### Using Volumes
You can use volumes options for external dependency, cache and etc.  
Set configurations in `.duci/config.yml`  

```yaml
volumes:
  - '/path/to/host/dir:/path/to/container/dir'
```

### Run Server
#### Add Environment Variable
This server needs environment variable `GITHUB_API_TOKEN` to create status.
```bash
export GITHUB_API_TOKEN=<your token>
```

## Server Settings
### Run Server
If you have already set $GOPATH, you can install it with the following command.
```bash
$ go get -u github.com/duck8823/duci
$ duci 
```

### Setting SSH
This server clone from github.com with **SSH** protocol
using private key `$HOME/.ssh/id_rsa` (default).  
Please set the public key of the pair at https://github.com/settings/keys.

### Server Configuration file
You can specify configuration file with `-c` option.
The configuration file must be yaml format.
Possible values ​​are as follows.

```yaml
server:
  workdir: '/path/to/tmp/duci'
  port: 8080
  sshKeyPath: '$HOME/.ssh/id_rsa'
  databasePath: '$HOME/.duci/db'
job:
  timeout: 600
  concurrency: `number of cpu`
```

You can check the default value.

```bash
$ duci -h
```

### Add Webhooks to GitHub repository
duci start to listen webhook with port `8080` and endpoint `/`.  
Add endpoint of duci to target repository.  
`https://github.com/<owner>/<repository>/settings/hooks`

## Using Docker
You can use Docker to run server.
```
$ docker run -p 8080:8080 \
             -e GITHUB_API_TOKEN=<your toekn> \
             -v /var/run/docker.sock:/var/run/docker.sock \
             -v ~/.ssh:/root/.ssh:ro \ 
             duck8823/duci
```

##### docker-compose for Windows
```bash
$ git clone https://github.com/duck8823/duci.git
$ cd duci
$ docker-compose -f docker-compose.win.yml up
```

##### docker-compose for Mac
```bash
$ git clone https://github.com/duck8823/duci.git
$ cd duci
$ docker-compose -f docker-compose.mac.yml up
```

## License
MIT License

Copyright (c) 2018 Shunsuke Maeda

See [LICENSE](./LICENSE) file
