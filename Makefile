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
