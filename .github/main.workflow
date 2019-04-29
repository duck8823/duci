workflow "main workflow" {
  on = "push"
  resolves = ["test", "check modified"]
}

action "lint" {
  uses = "docker://golangci/golangci-lint:latest"
  runs = ["golangci-lint", "run"]
  args = [
    "--disable-all",
    "--enable=gofmt",
    "--enable=vet",
    "--enable=gocyclo",
    "--enable=golint",
    "--enable=ineffassign",
    "--enable=misspell",
    "--deadline=5m"
  ]
}

action "download" {
  uses = "docker://golang:1.12"
  needs = ["lint"]
  env = {
    GOPATH = "/github/workspace/.go"
  }
  runs = "go"
  args = ["mod", "download"]
}

action "test" {
  uses = "docker://golang:1.12"
  needs = ["download"]
  env = {
    GOPATH = "/github/workspace/.go"
  }
  runs = "go"
  args = ["test", "./..."]
}

action "tidy" {
  uses = "docker://golang:1.12"
  needs = ["download"]
  env = {
    GOPATH = "/github/workspace/.go"
  }
  runs = "go"
  args = ["mod", "tidy"]
}

action "check modified" {
  uses = "docker://alpine/git:latest"
  needs = ["tidy"]
  runs = "sh"
  args = ["-c", "! git status | grep modified"]
}
