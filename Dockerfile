# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

RUN apk update && apk upgrade && apk --update add git make

EXPOSE 8080

WORKDIR /app

COPY . .

RUN go build -o server main.go

## Distribution
#FROM alpine:latest
#
#RUN apk update && apk upgrade && \
#    apk --update --no-cache add tzdata && \
#    mkdir /app
#
#WORKDIR /app
#
#EXPOSE 8080
#
#COPY --from=builder /app/engine /app
#
#CMD /app/engine