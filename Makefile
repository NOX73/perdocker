install: pull-images

docker-stop:
	sudo kill -QUIT `cat /var/run/docker.pid`

run: build run-perdocker

run-docker:
	sudo ./bin/docker -d
run-perdocker:
	./bin/perdocker

build:
	go build && mv perdocker ./bin/perdocker && chmod +x ./bin/perdocker

build-images: build-image-ruby build-image-nodejs build-image-go

build-image-ruby:
	docker build -rm -t="perdocker/ruby:attach" ./images/ruby/
build-image-nodejs:
	docker build -rm -t="perdocker/nodejs:attach" ./images/nodejs/
build-image-go:
	docker build -rm -t="perdocker/go:attach" ./images/go/

pull-images: pull-image-ruby pull-image-nodejs pull-image-go

pull-image-ruby:
	docker pull perdocker/ruby
pull-image-nodejs:
	docker pull perdocker/nodejs
pull-image-go:
	docker pull perdocker/go

