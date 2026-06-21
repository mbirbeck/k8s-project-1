FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o project-1-app

FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/project-1-app .
EXPOSE 8080
ENTRYPOINT ["/app/project-1-app"]

# GOOS=linux GOARCH=arm64