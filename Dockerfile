FROM golang:alpine
WORKDIR /app
COPY . /app
RUN go mod download
RUN go mod tidy
RUN go build -o cli ./cmd/cli/cli.go
ENTRYPOINT ["./cli"]
