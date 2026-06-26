FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o app-binary

FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/app-binary .
EXPOSE 8080
ENTRYPOINT ["/app/app-binary"]

# GOOS=linux GOARCH=arm64