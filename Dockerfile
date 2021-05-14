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
COPY static /app/static

ENTRYPOINT [ "./mqtt-listener-go", "start"]
