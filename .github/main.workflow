workflow "K8s-notify Testing" {
  resolves = ["GolangCI-Lint Action"]
  on = "pull_request"
}

action "GolangCI-Lint Action" {
  uses = "actions-contrib/golangci-lint@v0.1.0"
}
