FROM golang:latest AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o my_service cmd/app/main.go

FROM golang:latest
WORKDIR /app
COPY --from=builder /go/src/app/my_service /app/

CMD ["./my_service"]