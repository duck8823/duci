workflow "test" {
  on = "push"
  resolves = ["duci"]
}

action "duci" {
  uses = "./.duci/"
}
