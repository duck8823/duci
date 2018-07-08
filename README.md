# minimal-ci
[![Build Status](https://travis-ci.org/duck8823/minimal-ci.svg?branch=master)](https://travis-ci.org/duck8823/minimal-ci)
[![Coverage Status](https://coveralls.io/repos/github/duck8823/minimal-ci/badge.svg?branch=maintainance%2Fvarious)](https://coveralls.io/github/duck8823/minimal-ci?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/duck8823/minimal-ci)](https://goreportcard.com/report/github.com/duck8823/minimal-ci)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Minimal-ci is a small ci server.  
The job is triggered by pull request comment and push.  

## Target Repository
The target repository must have Dockerfile.  
I suggest to use `ENTRYPOINT` in Dockerfile.

e.g.
```Dockerfile
ENTRYPOINT ["mvn"]
CMD ["build"]
```

```Dockerfile
ENTRYPOINT ["fastlane"]
CMD ["build"]
```

When push to github, minimal-ci exec `mvn build` / `fastlane build`.  
And when comment `ci test` on github pull request, 
exec `mvn test` / `fastlane test` in docker container.  

## Run Server
### Add Environment Variable
This server needs environment variable `GITHUB_API_TOKEN` to create status.
```bash
export GITHUB_API_TOKEN=<your token>
```

### Setting SSH
This server clone from github.com with **SSH** protocol
using private key `$HOME/.ssh/id_rsa`.  
Please set the public key of the pair at https://github.com/settings/keys.

### Run Server
#### Locally
If you have already set $GOPATH, you can install it with the following command.
```bash
$ go get -u github.com/duck8823/minimal-ci
$ minimal-ci 
```

#### Using Docker
```
$ git clone https://github.com/duck8823/minimal-ci.git
$ docker build -t duck8823/minimal-ci .
$ docker run -e GITHUB_API_TOKEN=<your toekn> -v ~/.ssh:/root/.ssh:ro duck8823/minimal-ci
```

##### docker-compose for Windows
```bash
$ git clone https://github.com/duck8823/minimal-ci.git
$ cd minimal-ci
$ docker-compose -f docker-compose.win.yml up
```

##### docker-compose for Mac
```bash
$ git clone https://github.com/duck8823/minimal-ci.git
$ cd minimal-ci
$ docker-compose -f docker-compose.mac.yml up
```

#### Add Webhooks to GitHub repository
Add endpoint of minimal-ci to target repository.  
`https://github.com/<owner>/<repository>/settings/hooks`

