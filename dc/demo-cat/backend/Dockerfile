FROM golang:1.15-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /cat

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build main.go

FROM alpine:3.9

WORKDIR /app

COPY --from=build_base /cat/main /app/cat
COPY ./assets /app/assets

EXPOSE 7000

CMD ["./cat"]

