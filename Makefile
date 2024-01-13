GOLANG_CI_VERSION=v1.54.0

fmt:
	gofumpt -l -w .

test:
	go test ./... -v

lint/download:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.bin $(GOLANG_CI_VERSION) 
lint: lint/download
	./.bin/golangci-lint run -c ./.golangci.yml
lint/fix: lint/download
	./.bin/golangci-lint run -c ./.golangci.yml --fix


docker/build:
	docker build . -t ghcr.io/diezfx/idlegame-backend:latest -f "deployment/Dockerfile" --build-arg="APP_NAME=idlegame-backend"
docker/push: docker/build
	docker push ghcr.io/diezfx/idlegame-backend:latest

docker/up:
	docker compose up -d


migrate/drop:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path ./db/migrations drop




