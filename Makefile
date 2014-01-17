install: images_pull

run:
	./bin/perdocker

build:
	go build main.go && mv main ./bin/perdocker && chmod +x ./bin/perdocker

images_build: images_build_ruby images_build_nodejs

images_build_ruby:
	docker build -rm -t="perdocker/ruby" ./images/ruby/
images_build_nodejs:
	docker build -rm -t="perdocker/nodejs" ./images/nodejs/

images_pull: images_pull_ruby images_pull_nodejs

images_pull_ruby:
	docker pull perdocker/ruby
images_pull_nodejs:
	docker pull perdocker/nodejs

