build:
	go build
run:
	make build && ./plotbot-server
watch:
	nodemon --exec "make run" -e go && pkill plotbot-server
test:
	go test -v ./...
convey: 
	$$GOPATH/bin/goconvey -excludedDirs frontend,vendor,backups,uploads -port 9999
cover:
	make coverage
coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o=coverage.html && go tool cover -func=coverage.out
badge:
	gopherbadger -png=false -md=README.md
configure:
	make dep
dep:
	if ! [ -x "$(command -v dep)" ]; then\
	  go get github.com/golang/dep/cmd/dep;\
	fi && dep ensure;
