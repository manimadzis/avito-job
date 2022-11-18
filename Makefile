run: app
	./app -config ./config.yaml

go.mod:
	go mod init cmd/main.go
	go mod tidy

app: go.mod
	go build -o app cmd/main.go

docker:
	docker-compose run --build