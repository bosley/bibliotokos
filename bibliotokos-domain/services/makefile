all: build

build:
	go build -o bin/server ./cmd/server
	go build -o bin/admin ./cmd/admin
	cp .env bin/.env

server: build
	cd bin && ./server

admin: build
	./bin/admin

test:
	go test ./pkg/user/... -count=1 -v

clean:
	rm -f bin/server bin/admin bin/.env