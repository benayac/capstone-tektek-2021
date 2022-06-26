FROM golang:1.17-alpine AS builder
WORKDIR /workspace
COPY /web-service ./

# RUN go mod tidy

# RUN go mod download

# Build the API
ENV CGO_ENABLED = 0
ENV GOOS=linux
RUN go build -o /usr/local/bin/golang-docker .

# FINAL STAGE
FROM alpine:3.9
RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/bin/golang-docker /usr/local/bin/

RUN chown -R nobody:nogroup /usr/local/bin/golang-docker
USER nobody
EXPOSE 8080
CMD golang-docker