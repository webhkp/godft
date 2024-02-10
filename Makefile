build-windows:
	go build -o bin/godft.exe main.go 

build-linux:
	go build -o bin/godft main.go 

test:
	go test -v ./...

test-coverage:
	go test -cover ./...