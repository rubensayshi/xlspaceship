test: test-game test-client

test-game:
	go test -v github.com/rubensayshi/xlspaceship/pkg/ssgame

test-client:
	go test -v github.com/rubensayshi/xlspaceship/pkg/ssgame

