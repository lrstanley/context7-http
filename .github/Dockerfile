# syntax = docker/dockerfile:1.4

# backend
FROM golang:latest as build

COPY . /build
WORKDIR /build
RUN \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go \
	make build

# runtime
FROM alpine:3.23

RUN apk add --no-cache ca-certificates
COPY --from=build /build/context7-http /app/context7-http

EXPOSE 8080
WORKDIR /app
CMD ["/app/context7-http"]
