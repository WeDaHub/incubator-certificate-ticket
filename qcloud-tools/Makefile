build:
	GOOS=linux GOARCH=amd64 go build -o bin/cert-monitor cmd/main.go
	chmod +x bin/cert-monitor

clean:
	rm -rf bin/*
	go clean
