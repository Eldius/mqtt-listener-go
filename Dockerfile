FROM node:lts as nodebuilder
#FROM node:lts as nodebuilder
#FROM node:14.16.1-alpine3.13 as nodebuilder
#FROM node:lts-buster as nodebuilder
#FROM node:14.16.1-buster-slim as nodebuilder

WORKDIR /app
COPY ./static /app

ENV REACT_APP_BACKEND_ENDPOINT="http://192.168.100.195/mqtt-listener"

RUN yarn install
RUN yarn build

FROM golang:1.16-alpine3.13 as gobuilder

WORKDIR /app
COPY . /app

RUN apk add --no-cache git make build-base
RUN go build -v -a -ldflags '-extldflags "-static"' .
RUN chmod +x /app/mqtt-listener-go

FROM alpine:3.13

EXPOSE 8080

WORKDIR /app

COPY --chown=0:0 --from=gobuilder /app/mqtt-listener-go /app
COPY --chown=0:0 --from=nodebuilder /app/build /app/static

ENTRYPOINT [ "./mqtt-listener-go", "start"]
