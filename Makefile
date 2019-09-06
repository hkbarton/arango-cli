PROJECT := arango-cli

all: build

build:
	go build -o build/${PROJECT}

clean:
	rm -rf ./build

run:
	make clean
	make
	./build/${PROJECT} -h 127.0.0.1 -p 8529