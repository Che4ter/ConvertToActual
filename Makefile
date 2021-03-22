LDFLAGS=-ldflags "-s -w -X main.Version=$(shell git describe --abbrev=0  --always --tags)"

compile:
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/convertToActual ${LDFLAGS} ./main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/convertToActual.exe ${LDFLAGS} ./main.go
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o bin/convertToActual-darwin ${LDFLAGS} ./main.go
clean:
	rm -f *_converted.csv bin/*
