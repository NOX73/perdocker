install: images_pull

run:
	go run main.go

images_build: images_build_ruby

images_build_ruby:
	docker build -rm -t="perdocker/ruby" ./images/ruby/

images_pull: images_pull_ruby

images_pull_ruby:
	docker pull perdocker/ruby

