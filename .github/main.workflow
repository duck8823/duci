workflow "main workflow" {
  on = "push"
  resolves = ["test", "check modified", "goreportcard"]
}

action "download" {
  uses = "docker://golang:1.11"
  env = {
    GOPATH = "/github/workspace/.go"
  }
  runs = "go"
  args = ["mod", "download"]
}

action "test" {
  uses = "docker://golang:1.11"
  needs = ["download"]
  env = {
    GOPATH = "/github/workspace/.go"
  }
  runs = "go"
  args = ["test", "./..."]
}

action "tidy" {
  uses = "docker://golang:1.11"
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

action "goreportcard" {
  uses = "docker://duck8823/goreportcard:latest"
  args = ["-t", "100"]
}