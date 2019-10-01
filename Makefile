build:
	go build -o anticors

run: build
	./anticors

deploy: build
	rsync -rvz anticors root@toorop.fr:/var/www/anticors.toorop.fr/