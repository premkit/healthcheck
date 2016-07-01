.PHONY: clean godep deps run test build all

clean:
	rm -f ./bin/healthcheck

run:
	./bin/healthcheck daemon

test:
	govendor test +local

build:
	mkdir -p ./bin
	govendor build -o ./bin/healthcheck .

shell:
	docker run --rm -it -P --name healthcheck \
                -p 8201:80 \
                -e PREMKIT_ROUTER=$(PREMKIT_ROUTER) \
                -e ADVERTISE_ADDRESS="http://$(shell ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+'):8201" \
		-v `pwd`:/go/src/github.com/premkit/healthcheck \
                -v `pwd`/data:/data \
		premkit/healthcheck:dev

all: build test
