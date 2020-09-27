FROM golang:1.15 AS build-env

ADD . /bload
WORKDIR /bload

RUN go build -o bload.out

FROM debian:buster

EXPOSE 3000

WORKDIR /
COPY --from=build-env /bload/static/ /static
COPY --from=build-env /bload/views /views
COPY --from=build-env /bload/bload.out /
