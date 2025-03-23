clean:
	rm -f app
build:
	cd cmd/golang-server-layout && go build -o ../../app
run: build
	./app

