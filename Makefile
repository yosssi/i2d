build:
	rm -rf ./artifacts
	mkdir ./artifacts
	GOOS=linux   GOARCH=amd64 go build -o ./artifacts/i2d-linux-amd64       cmd/i2d/main.go
	GOOS=linux   GOARCH=386   go build -o ./artifacts/i2d-linux-386         cmd/i2d/main.go
	GOOS=darwin  GOARCH=amd64 go build -o ./artifacts/i2d-darwin-amd64      cmd/i2d/main.go
	GOOS=darwin  GOARCH=386   go build -o ./artifacts/i2d-darwin-386        cmd/i2d/main.go
	GOOS=windows GOARCH=amd64 go build -o ./artifacts/i2d-windows-amd64.exe cmd/i2d/main.go
	GOOS=windows GOARCH=386   go build -o ./artifacts/i2d-windows-386.exe   cmd/i2d/main.go
