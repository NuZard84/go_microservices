FROM golang:1.23.8-alpine3.18 AS build

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /go/src/github.com/NuZard84/go_microservices

COPY go.mod go.sum ./

COPY vendor vendor

COPY order order

RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

FROM alpine:3.18

WORKDIR /usr/bin

COPY --from=build /go/bin .

EXPOSE 8080

CMD [ "./app" ]


