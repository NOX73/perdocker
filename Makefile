images: images_ruby

images_ruby:
	docker build -rm -t="ruby:2.1.0" ./images/ruby/

run:
	go run main.go
