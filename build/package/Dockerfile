FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/server -a -o main .

FROM alpine:3.20  
WORKDIR /root/
COPY --from=builder /app/cmd/server/main .
EXPOSE 8080
CMD ["./main"] 