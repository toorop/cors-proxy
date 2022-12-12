build:
	go build -o ovhcors

run: build
	./ovhcors

deploy: build
	rsync -rvz ovhcors root@toorop.fr:/var/www/cors.dpp.st/