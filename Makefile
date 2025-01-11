run:
	@go run -race . -v

dbuild:
	@docker build -t dmdhrumilmistry/defect-detect-backend --progress text -f backend.Dockerfile . 

scan-vulns:
	@trivy image dmdhrumilmistry/defect-detect-backend

docker-local: dbuild scan-vulns

build:
	@go build -o bin/defect-detect .

test:
	@go test -cover -v ./...

bump:
	@go get -u ./...
	@go mod tidy

local: docker scan-vulns