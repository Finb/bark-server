workflow "build and push to dockerhub" {
  on = "push"
  resolves = ["login", "test", build", "push"]
}

action "login" {
  uses = "actions/docker/login@8cdf801b322af5f369e00d85e9cf3a7122f49108"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "test" {
  uses="cedrickring/golang-action@1.2.0"
}

action "build" {
  needs = ["test"]
  uses = "actions/docker/cli@master"
  args = "build -t metrue/bark-server:latest ."
}

action "push" {
  needs = ["build", "login"]
  uses = "actions/docker/cli@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
  args = "push metrue/bark-server:latest"
}
