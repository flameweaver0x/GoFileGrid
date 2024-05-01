FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gofilegrid server.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gofilegrid .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./gofilegrid"]