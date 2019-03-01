workflow "main workflow" {
  on = "push"
  resolves = ["test", "check modified"]
}

action "lint" {
  uses = "docker://duck8823/gometalinter:latest"
  args = [
    "--disable-all",
    "--enable=gofmt",
    "--enable=vet",
    "--enable=gocyclo", "--cyclo-over=15",
    "--enable=golint", "--min-confidence=0.85", "--vendor",
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
