FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go build -o project-1-app

FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/project-1-app .
EXPOSE 8080
CMD [ "sh", "-c", "/app/project-1-app"]
