test: test-game test-client

test-game:
	go test -v github.com/rubensayshi/xlspaceship/pkg/ssgame

test-client:
	go test -v github.com/rubensayshi/xlspaceship/pkg/ssclient

coverage:
	go test -v github.com/rubensayshi/xlspaceship/pkg/ssgame -cover -coverprofile=coverage1.out
	go test -v github.com/rubensayshi/xlspaceship/pkg/ssclient -cover -coverprofile=coverage2.out
	go run vendor/github.com/wadey/gocovmerge/gocovmerge.go coverage1.out coverage2.out > coverage.out
	go tool cover -func=coverage.out

build-gui:
	cd ./gui && bower update && gulp

watch-build-gui:
	cd ./gui && bower update && gulp watch --live-reload

build-gui-statik:
	rm ./statik/statik.go || echo ""
	go run vendor/github.com/rakyll/statik/statik.go -src=./gui/web/www/ -f -m -Z

build: build-gui build-gui-statik build-linux build-windows

build-windows:
	 GOOS=windows GOARCH=amd64 go build -tags statik -o bin/xlspaceship-win64.exe main.go
	 GOOS=windows GOARCH=386 go build -tags statik -o bin/xlspaceship-win386.exe main.go

build-linux:
	 GOOS=linux GOARCH=amd64 go build -tags statik -o bin/xlspaceship-linux64 main.go
	 GOOS=linux GOARCH=386 go build -tags statik -o bin/xlspaceship-linux386 main.go
