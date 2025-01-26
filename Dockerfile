FROM golang:1.23
WORKDIR /api

# dependencies
COPY go.mod ./
RUN go mod download

# compile
COPY main.go ./
COPY lib lib
RUN mkdir -p dist &&  CGO_ENABLED=0 GOOS=linux go build -o dist/asset-pricing

# entrypoint
ENTRYPOINT [ "dist/asset-pricing" ]