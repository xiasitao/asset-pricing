FROM golang:1.23 AS build
WORKDIR /dist

# dependencies
COPY go.mod ./
RUN go mod download

# compile
COPY main.go ./
COPY lib lib
RUN CGO_ENABLED=0 GOOS=linux go build -o app.bin

# entrypoint
FROM golang:1.23 AS run
WORKDIR /app
COPY --from=build /dist/app.bin .
ENTRYPOINT [ "./app.bin" ]