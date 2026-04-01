FROM golang:latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ledger-core ./cmd/server/main.go

EXPOSE 8080
CMD ["./ledger-core"]