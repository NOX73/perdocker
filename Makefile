images: images_ruby

images_ruby:
	docker build -t="perdocker:ruby" ./images/ruby/

run:
	go run main.go
