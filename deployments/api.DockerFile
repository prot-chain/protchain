FROM --platform=linux/amd64 golang:1.21.1-alpine AS builder
WORKDIR /app
ADD go.mod go.sum ./

RUN go mod download

ADD . .
RUN go build -o main ./cmd/protchainapi/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
