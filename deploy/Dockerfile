FROM golang:1.14.7-alpine3.12 AS build

WORKDIR /src/

COPY go.mod go.sum ./
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN : \
    && CGO_ENABLED=0 go build ./cmd/templatemap \
    ;


FROM alpine:3.12

WORKDIR /root

COPY --from=build /src/templatemap /root/

ENTRYPOINT ["/root/templatemap"]
