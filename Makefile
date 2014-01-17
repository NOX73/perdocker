install: pull-images

run:
	./bin/docker -d & ./bin/perdocker

run-docker:
	./bin/docker -d
run-perdocker:
	./bin/perdocker

build:
	go build main.go && mv main ./bin/perdocker && chmod +x ./bin/perdocker

build-images: build-image-ruby build-image-nodejs

build-image-ruby:
	docker build -rm -t="perdocker/ruby" ./images/ruby/
build-image-nodejs:
	docker build -rm -t="perdocker/nodejs" ./images/nodejs/

pull-images: pull-image-ruby pull-image-nodejs

pull-image-ruby:
	docker pull perdocker/ruby
pull-image-nodejs:
	docker pull perdocker/nodejs

