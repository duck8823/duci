workflow "main workflow" {
  on = "push"
  resolves = ["test"]
}

action "test" {
  uses = "docker://golang:1.11"
  runs = "go"
  args = ["test", "./..."]
}
