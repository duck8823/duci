# minimal-ci
Minimal-ci is a small ci server.  
The job is triggered by pull request comment.  

## Target Repository
The target repository must have Dockerfile.  
I suggest to use `ENTRYPOINT` in Dockerfile.

e.g.
```Dockerfile
ENTRYPOINT ["mvn"]
```

```Dockerfile
ENTRYPOINT ["fastlane"]
```

When comment `ci test` on github pull request, 
minimal-ci exec `mvn test` / `fastlane test` in docker container.  

## Run Server
### Add Environment Variable
This server needs environment variable `GITHUB_API_TOKEN` to create status.
```bash
export GITHUB_API_TOKEN=<your token>
```

### Run Server
#### Locally
If you have already set $GOPATH, you can install it with the following command.
```bash
$ go get -u github.com/duck8823/minimal-ci
$ minimal-ci 
```

#### Using Docker
##### Docker for Windows
```bash
$ git clone https://github.com/duck8823/minimal-ci.git
$ cd minimal-ci
$ docker-compose -f docker-compose.win.yml up
```

##### Docker for Mac
```bash
$ git clone https://github.com/duck8823/minimal-ci.git
$ cd minimal-ci
$ docker-compose -f docker-compose.mac.yml up
```

#### Add Webhooks to GitHub repository
Add endpoint of minimal-ci to target repository.  
`https://github.com/<owner>/<repository>/settings/hooks`