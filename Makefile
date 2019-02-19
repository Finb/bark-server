build:
	go build -o bin/bark bark/bark.go
start: build
	./bin/bark -ip=0.0.0.0 -port=80
